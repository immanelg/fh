package main

import (
	"encoding/json"
	"net/http"
	"os"
)

type createDirResp struct {
	Entry fileMeta
}

func apiCreateDir(path string, w http.ResponseWriter, r *http.Request) {
	// TODO: check that parent is a directory and that there's no files at this path
    err := os.MkdirAll(path, 0o777)
	if err != nil {
		internalError(w, err)
		return
	}

	var respModel createFileResp
	entry, err := fileMetaOf(path)
	if err != nil {
		internalError(w, err)
		return
	}
	respModel.Entry = entry
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	e := json.NewEncoder(w)
	e.Encode(respModel)
}
