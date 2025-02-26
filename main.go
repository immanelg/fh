package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var dir string = "."
var host string = "localhost"
var port string = "8080"

func setupConfig() {
	flag.StringVar(&dir, "d", ".", "directory to serve")
	flag.StringVar(&port, "p", "8080", "port")
	flag.Parse()

	if fi, err := os.Stat(dir); err != nil || !fi.Mode().IsDir() {
		log.Fatal("bad directory 😠\n")
	}

}

func main() {
    setupConfig()

	s := http.NewServeMux()
	s.HandleFunc("POST /list-dir", apiListDir)
	s.HandleFunc("POST /read-file", apiReadFile)
	s.HandleFunc("POST /create-file", apiCreateFile)
	s.HandleFunc("POST /create-dir", apiCreateDir)
	s.HandleFunc("POST /delete", apiDelete)
	s.HandleFunc("POST /copy", apiCopy)
	s.HandleFunc("POST /metadata", apiMetadata)

	log.Printf("starting on http://%s:%s\n", host, port)
	log.Printf("serving dir %s\n", dir)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), s))
}
