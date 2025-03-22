package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var basedir string = "."
var host string = "localhost"
var port string = "8080"

func configure() {
	flag.StringVar(&basedir, "d", ".", "directory to serve")
	flag.StringVar(&host, "H", "localhost", "host")
	flag.StringVar(&port, "p", "8080", "port")
	flag.Parse()

	if fi, err := os.Stat(basedir); err != nil || !fi.Mode().IsDir() {
		log.Fatalf("bad directory 😠 %s\n", err)
	}
    
    // ?
    // if err := os.Chdir(basedir); err != nil {
    //     log.Fatalf("can't Chdir: %s\n", err)
    // }
}

func handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.String())
    q := r.URL.Query()
	headers := r.Header

	path := filepath.Join(basedir, r.URL.Path)

	switch r.Method {
	case http.MethodGet:
        m := q.Get("metadata")
        if m == "1" {
            apiMetadata(path, w, r)
        } else {
            apiReadFile(path, w, r)
        }
	case http.MethodPost:
        // if moveFrom, exists := headers["Move-From"]; exists {
            // apiMove(w, r)
            // break
        // }

        if src := headers.Get("Source-Path"); src != "" {
            src = filepath.Join(basedir, src)
            op := headers.Get("Operation")
            if op == "Copy" {
                apiCopy(path, src, w, r)
            } else if op == "Move" {

            }
            break
        }

        t := q.Get("resource-type") 
        if t == "dir" {
            apiCreateDir(path, w, r)
        } else {
            apiCreate(path, w, r)
        }
	case http.MethodDelete:
        apiDelete(path, w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	configure()

	s := &http.Server{
		Addr:           net.JoinHostPort(host, port),
		Handler:        http.HandlerFunc(handle),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	log.Printf("serving dir %s\n", basedir)
    log.Fatal(s.ListenAndServe())

    // go log.Fatal(s.ListenAndServe())
    // c := make(chan os.Signal, 1)
    // signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
    // ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    // defer cancel()
    // log.Fatal(s.Shutdown(ctx))
}
