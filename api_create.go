package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type createFileResp struct {
	Entry fileMeta
}

func apiCreate(path string, w http.ResponseWriter, r *http.Request) {
    // TODO: check that parent is a directory
    f, err := os.Create(path)
    if err != nil {
        log.Printf("error: %s\n", err.Error())
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    defer f.Close()

    _, err = io.Copy(f, r.Body)
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
