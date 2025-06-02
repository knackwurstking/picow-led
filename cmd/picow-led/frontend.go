//go:build frontend
// +build frontend

package main

import (
	"embed"
	"io/fs"
)

//go:embed frontend-build
var frontendBuild embed.FS

func frontend() (f fs.FS, ok bool) {
	f, err := fs.Sub(frontendBuild, "frontend-build")
	if err != nil {
		panic(err)
	}

	return f, true
}
