package main

import "time"

type fileEntry struct {
	Name    string
	Path    string
	Type    string
	Size    uint64
	ModTime time.Time
    // TODO: unix permission bits?
}
