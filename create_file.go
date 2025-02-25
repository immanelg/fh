package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type createFileReq struct {
	Path string
	Type string
	// TOOD: more metadata
}

type createFileResp struct {
	Entry fileEntry
}

func apiCreateFile(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "Invalid multipart form data", http.StatusBadRequest)
		return
	}

	var reqModel createFileReq
	var fileFound bool
	var payloadFound bool
	var path string
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Error reading form part", http.StatusInternalServerError)
			return
		}

		if part.FormName() == "payload" {
			payloadFound = true
			d := json.NewDecoder(part)
			err := d.Decode(&reqModel)
			if err != nil {
				http.Error(w, "Invalid payload", http.StatusBadRequest)
				return
			}
		} else if part.FormName() == "file" {
			fileFound = true

			if reqModel.Type != "File" {
				http.Error(w, "Unexpected file Type", http.StatusBadRequest)
				return
			}
			if !payloadFound {
				http.Error(w, "File must come after payload in form data", http.StatusBadRequest)
				return
			}

			path = filepath.Join(dir, reqModel.Path)

			// TODO: check that parent is a directory
			// TODO: security
			f, err := os.Create(path)
			if err != nil {
				log.Printf("error: %s\n", err.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			defer f.Close()

			_, err = io.Copy(f, part)
			if err != nil {
				internalError(w, err)
				return
			}
		} else {
			http.Error(w, "Unexpected form value", http.StatusBadRequest)
			return
		}

		part.Close()
	}

	if !payloadFound {
		http.Error(w, "Expected payload in form data", http.StatusBadRequest)
		return
	}

	if !fileFound {
		http.Error(w, "Expected file in form data", http.StatusBadRequest)
		return
	}

	var respModel createFileResp
	fi, err := os.Stat(path)
	if err != nil {
		internalError(w, err)
		return
	}
	respModel.Entry.Name = fi.Name()
	respModel.Entry.Path = path
	mode := fi.Mode()
	if mode.IsDir() {
		respModel.Entry.Type = "Dir"
	} else {
		respModel.Entry.Type = "File"
	}
	respModel.Entry.Size = uint64(fi.Size())
	respModel.Entry.ModTime = fi.ModTime()

	w.WriteHeader(http.StatusCreated)
	e := json.NewEncoder(w)
	e.Encode(respModel)
}
