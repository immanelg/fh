package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type listReq struct {
	Path string
	// TODO: options
}
type listResp struct {
	Entries []fileEntry
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
	}
	var respModel listResp
	respModel.Entries = make([]fileEntry, len(entries))
	for i, entry := range entries {
		respModel.Entries[i] = fileEntry{
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
