package ui

import (
	"embed"
	"io/fs"
)

//go:embed dist
var dist embed.FS

func Dist() fs.FS {
	data, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}

	return data
}
