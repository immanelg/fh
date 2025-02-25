package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type createDirReq struct {
	Path string
}

type createDirResp struct {
	Entry fileEntry
}

func apiCreateDir(w http.ResponseWriter, r *http.Request) {
	var reqModel createDirReq

	d := json.NewDecoder(r.Body)
	err := d.Decode(&reqModel)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// TODO: check that parent is a directory and that there's no files at this path
	// TODO: security
	err = os.MkdirAll(filepath.Join(dir, reqModel.Path), 0o777)
	if err != nil {
		internalError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	// TODO: write metadata json
}
