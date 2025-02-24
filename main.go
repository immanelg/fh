package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var dir string = "."
var host string = "localhost"
var port string = "8080"

type listReq struct {
	Path string
	// TODO: options
}
type dirEntry struct {
	Name    string
	Path    string
	Type    string
	Size    uint64
	ModTime time.Time
}
type listResp struct {
	Entries []dirEntry
}
type readReq struct {
	Path string
}

func apiListDir(w http.ResponseWriter, r *http.Request) {
	var reqModel listReq
	d := json.NewDecoder(r.Body)
	err := d.Decode(&reqModel)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	path := filepath.Join(dir, reqModel.Path)
	fi, err := os.Stat(path)
	if handleIoError(w, err) {
		return
	}
	if !fi.IsDir() {
		// TODO: write info for one file
	}

	// TODO: security
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	var respModel listResp
	respModel.Entries = make([]dirEntry, len(entries))
	for i, entry := range entries {
		respModel.Entries[i] = dirEntry{
			Name: entry.Name(),
			Path: filepath.Join(path, entry.Name()),
		}
		if fi, err := entry.Info(); err == nil {
			mode := fi.Mode()
			// TODO: symlinks https://pkg.go.dev/io/fs#FileMode
			if mode.IsDir() {
				respModel.Entries[i].Type = "Dir"
			} else {
				respModel.Entries[i].Type = "File"
			}
			respModel.Entries[i].Size = uint64(fi.Size())
			respModel.Entries[i].ModTime = fi.ModTime()
		}
	}

	w.Header().Add("Content-Type", "application/json")
	e := json.NewEncoder(w)
	e.Encode(respModel)
}

func handleIoError(w http.ResponseWriter, err error) (haveError bool) {
	if err == nil {
		return false
	}
	if os.IsNotExist(err) {
		http.Error(w, "Not found", http.StatusNotFound)
		return true
	} else {
		log.Printf("error: %s\n", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return true
	}
}
func apiRead(w http.ResponseWriter, r *http.Request) {
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
	if handleIoError(w, err) {
		return
	}
	if !fi.Mode().IsRegular() {
		// TODO: links
		http.Error(w, "Forbidden: not a regular file", http.StatusForbidden)
		return
	}

	f, err := os.Open(path)
	if handleIoError(w, err) {
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
}

type createReq struct {
	Path string
	Type string
	// TOOD: more metadata
}

func apiCreate(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "Invalid multipart form data", http.StatusBadRequest)
		return
	}
    
	var reqModel createReq
	var fileFound bool
    var payloadFound bool
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
				http.Error(w, "Found unexpected file in form data", http.StatusBadRequest)
				return
			}
			if !payloadFound {
				http.Error(w, "File must come after payload in form data", http.StatusBadRequest)
				return
			}

            path := filepath.Join(dir, reqModel.Path)

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
				log.Printf("error: %s\n", err.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
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

	if reqModel.Type == "File" {
        if !fileFound {
            http.Error(w, "Expected file in form data", http.StatusBadRequest)
            return
        }
        return // we are done already because we streamed the file immediatly above
	} else if reqModel.Type == "Dir" {
		// TODO: check that parent is a directory and that there's no files at this path
		// TODO: security
		err := os.MkdirAll(filepath.Join(dir, reqModel.Path), 0o777)
		if err != nil {
			log.Printf("error: %s\n", err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
        http.Error(w, "Invalid Type in payload", http.StatusBadRequest)
        return
    }
}

type copyReq struct {
	Src string
    Dst string
    // TODO: Symlink, Copy, Move
}

func apiCopy(w http.ResponseWriter, r *http.Request) {
	var reqModel copyReq
	d := json.NewDecoder(r.Body)
	err := d.Decode(&reqModel)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

    srcpath := filepath.Join(dir, reqModel.Src)
    dstpath := filepath.Join(dir, reqModel.Dst)
    log.Println(srcpath, dstpath)

    // stat
    srcfi, err := os.Stat(srcpath)
    if err != nil {
		log.Printf("error: %s\n", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    if !srcfi.Mode().IsRegular() {
		http.Error(w, "Src is not a regular file", http.StatusBadRequest)
		return
    }
    dstfi, err := os.Stat(dstpath)
    if !os.IsNotExist(err) && err != nil {
		log.Printf("error: %s\n", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
	if err == nil /* file exists */ && dstfi.IsDir() {
		dstpath = filepath.Join(dstpath, filepath.Base(srcpath))
	}

    // open and create
    srcf, err := os.Open(srcpath)
    if err != nil {
		log.Printf("error: %s\n", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    defer srcf.Close()
    dstf, err := os.Create(dstpath)
    if err != nil {
		log.Printf("error: %s\n", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    defer dstf.Close()

    // copy
    _, err = dstf.ReadFrom(srcf)
    if err != nil {
		log.Printf("error: %s\n", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    err = dstf.Sync()
    if err != nil {
		log.Printf("error: %s\n", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
}

func main() {
    flag.StringVar(&dir, "d", ".", "directory to serve")
    flag.StringVar(&port, "p", "8080", "port")
    flag.Parse()

    if fi, err := os.Stat(dir); err != nil || !fi.Mode().IsDir() {
        log.Fatal("bad directory 😠\n")
    }

	s := http.NewServeMux()
	s.HandleFunc("POST /list", apiListDir)
	s.HandleFunc("POST /read", apiRead)
	s.HandleFunc("POST /create", apiCreate)
	s.HandleFunc("POST /copy", apiCopy)

	log.Printf("starting on http://%s:%s\n", host, port)
	log.Printf("serving dir %s\n", dir)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), s))
}
