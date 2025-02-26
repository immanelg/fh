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
	Entry fileMeta
}

func apiCreateDir(w http.ResponseWriter, r *http.Request) {
	var reqModel createDirReq

	d := json.NewDecoder(r.Body)
	err := d.Decode(&reqModel)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// TODO: security
    path := filepath.Join(dir, reqModel.Path)

	// TODO: check that parent is a directory and that there's no files at this path
	err = os.MkdirAll(path, 0o777)
	if err != nil {
		internalError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	var respModel createFileResp
    entry, err := fileMetaOf(path)
	if err != nil {
		internalError(w, err)
		return
	}
	respModel.Entry = entry
	w.WriteHeader(http.StatusCreated)
	e := json.NewEncoder(w)
	e.Encode(respModel)

	// TODO: write metadata json
}
