package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type copyResp struct {
	Entry fileMeta
}

func apiCopy(dstpath string, srcpath string, w http.ResponseWriter, r *http.Request) {
	srcfi, err := os.Stat(srcpath)
	if err != nil {
		internalError(w, err)
		return
	}
	if !srcfi.Mode().IsRegular() {
		http.Error(w, "source is not a regular file", http.StatusBadRequest)
		return
	}
	dstfi, err := os.Stat(dstpath)
	if !os.IsNotExist(err) && err != nil {
		internalError(w, err)
		return
	}
	if err == nil /* file exists */ && dstfi.IsDir() {
		dstpath = filepath.Join(dstpath, filepath.Base(srcpath))
	}

	// open and create
	srcf, err := os.Open(srcpath)
	if err != nil {
		internalError(w, err)
		return
	}
	defer srcf.Close()
	dstf, err := os.Create(dstpath)
	if err != nil {
		internalError(w, err)
		return
	}
	defer dstf.Close()

	// copy
	_, err = dstf.ReadFrom(srcf)
	if err != nil {
		internalError(w, err)
		return
	}
	err = dstf.Sync()
	if err != nil {
		internalError(w, err)
		return
	}

	var respModel copyResp
	entry, err := fileMetaOf(dstpath)
	if err != nil {
		internalError(w, err)
		return
	}
	respModel.Entry = entry
	w.WriteHeader(http.StatusCreated)
	e := json.NewEncoder(w)
	e.Encode(respModel)
}
