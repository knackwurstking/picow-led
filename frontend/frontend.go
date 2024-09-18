package frontend

import (
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
)

//go:embed dist
var dist embed.FS

func GetFS() fs.FS {
	fsys, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(fmt.Sprintf("Get sub fs: %s", err.Error()))
	}
	return fsys
}
