package main

import (
	"log"
	"net/http"
	"os"
)

type deleteReq struct {
	Path string
}

func apiDelete(path string, w http.ResponseWriter, r *http.Request) {
    err := os.RemoveAll(path)
	// NOTE: If the path does not exist, RemoveAll returns nil
	if err != nil {
		log.Printf("error: %s\n", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
