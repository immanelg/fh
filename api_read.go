package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func apiReadFile(path string, w http.ResponseWriter, r *http.Request) {
	fi, err := os.Stat(path)
	if err != nil {
		handleNotFoundOrInternalErr(w, err)
		return
	}
	if !fi.Mode().IsRegular() {
		// TODO: links
		http.Error(w, "Not a regular file", http.StatusForbidden)
		return
	}

	if fi.IsDir() {
	}

	if mt := r.Header.Get("If-Modified-Since"); mt != "" {
		if mt, err := time.Parse(http.TimeFormat, mt); err == nil && fi.ModTime().After(mt) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	f, err := os.Open(path)
	if err != nil {
		handleNotFoundOrInternalErr(w, err)
		return
	}
	defer f.Close()

	w.Header().Set("Last-Modified", fi.ModTime().Format(http.TimeFormat))

	if _, err = io.Copy(w, f); err != nil {
		log.Printf("error: %s", err.Error())
	}
	// TODO: metadata?
}
