package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type deleteReq struct {
	Path string
}

func apiDelete(w http.ResponseWriter, r *http.Request) {
	var reqModel deleteReq
	d := json.NewDecoder(r.Body)
	err := d.Decode(&reqModel)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	path := filepath.Join(dir, reqModel.Path)

	err = os.RemoveAll(path)
	// NOTE: If the path does not exist, RemoveAll returns nil
	if err != nil {
		log.Printf("error: %s\n", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
