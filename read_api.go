package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type readReq struct {
	Path string
}

func apiReadFile(w http.ResponseWriter, r *http.Request) {
	var reqModel readReq
	d := json.NewDecoder(r.Body)
	err := d.Decode(&reqModel)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// TODO(security): Join with absolute path
	path := filepath.Join(dir, reqModel.Path)

	fi, err := os.Stat(path)
	if err != nil {
		handleNotFoundOrInternalErr(w, err)
		return
	}
	if !fi.Mode().IsRegular() {
		// TODO: links
		http.Error(w, "Forbidden: not a regular file", http.StatusForbidden)
		return
	}

	f, err := os.Open(path)
	if err != nil {
		handleNotFoundOrInternalErr(w, err)
		return
	}
	defer f.Close()

	if fi.IsDir() {
	}

	if mt := r.Header.Get("If-Modified-Since"); mt != "" {
		if mt, err := time.Parse(http.TimeFormat, mt); err == nil && fi.ModTime().After(mt) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}
	w.Header().Set("Last-Modified", fi.ModTime().Format(http.TimeFormat))

	if _, err = io.Copy(w, f); err != nil {
		log.Printf("error: %s", err.Error())
	}
	// TODO: multipart for metadata + file?
}
