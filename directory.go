package main

import (
	"path/filepath"
	"strings"
)

func localizePath(path string) string {
	return strings.TrimPrefix(basedir, path)
}

func join(parent, child string) {
	filepath.Join(parent, child)
}
