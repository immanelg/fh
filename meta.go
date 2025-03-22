package main

import (
	"os"
	"time"
)

type fileMeta struct {
	Name    string
	Path    string
	Type    string
	Size    uint64
	ModTime time.Time
	// TODO: unix permission bits?
}

func fileMetaOf(path string) (fileMeta, error) {
	var e fileMeta
	fi, err := os.Stat(path)
	if err != nil {
		return e, err
	}
	e.Name = fi.Name()
	// FIXME: here and in list-dir, strip root directory prefix from path.
	e.Path = path
	mode := fi.Mode()
	if mode.IsDir() {
		e.Type = "Dir"
	} else {
		e.Type = "File"
	}
	e.Size = uint64(fi.Size())
	e.ModTime = fi.ModTime()
	return e, nil
}
