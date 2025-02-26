package main

import (
	"log"
	"net/http"
	"os"
)

func internalError(w http.ResponseWriter, err error) {
	log.Printf("error: %s\n", err.Error())
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func handleNotFoundOrInternalErr(w http.ResponseWriter, err error) bool {
	if os.IsNotExist(err) {
		http.Error(w, "Not found", http.StatusNotFound)
		return true
	} else {
		internalError(w, err)
		return false
	}
}
