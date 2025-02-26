package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type listReq struct {
	Path string
	// TODO: options
}
type listResp struct {
	Entries []fileMeta
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
	if err != nil {
		handleNotFoundOrInternalErr(w, err)
		return
	}
	if !fi.IsDir() {
		// TODO: write info for one file
	}

	// TODO: security
	entries, err := os.ReadDir(path)
	if err != nil {
		handleNotFoundOrInternalErr(w, err)
        return
	}
	var respModel listResp
	respModel.Entries = make([]fileMeta, len(entries)) // preallocate
	for i, entry := range entries {
        meta, err := fileMetaOf(filepath.Join(path, entry.Name()))
        if err != nil {
            log.Printf("error when listing dir: %v", err.Error())
        }
        respModel.Entries[i] = meta
	}

	w.Header().Add("Content-Type", "application/json")
	e := json.NewEncoder(w)
	e.Encode(respModel)
}
