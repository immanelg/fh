package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type copyReq struct {
	// TODO: multiple Src to one dir
	Src  string
	Dst  string
	Kind string // TODO: Copy, Move, SymlinkLink
}

type copyResp struct {
	Entry fileEntry
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

	// stat
	srcfi, err := os.Stat(srcpath)
	if err != nil {
		internalError(w, err)
		return
	}
	if !srcfi.Mode().IsRegular() {
		http.Error(w, "Src is not a regular file", http.StatusBadRequest)
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
	// TODO: write metadata response
}
