package web

import (
	"embed"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type (
	File  string
	Key   string
	Value string
)

type ServeData struct {
	FS               fs.FS
	ServerPathPrefix string
	TemplateData     any
}

func StripSubFromFS(t embed.FS) fs.FS {
	f, err := fs.Sub(t, "templates")
	if err != nil {
		panic(err)
	}
	return f
}

func ServeGET(e *echo.Echo, header map[File]map[Key][]Value, data ServeData) {
	for file, headers := range header {
		e.GET(data.ServerPathPrefix+string(file), func(c echo.Context) error {
			t, err := template.ParseFS(data.FS, string(file))
			if err != nil {
				return err
			}

			for key, values := range headers {
				for _, value := range values {
					c.Response().Header().Add(string(key), string(value))
				}
			}

			return t.Execute(c.Response(), data.TemplateData)
		})
	}
}

func GenerateFile(tFS fs.FS, outPath string, filePath string, data any) error {
	err := os.MkdirAll(outPath, 0700)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(outPath, filePath))
	if err != nil {
		return err
	}

	t, err := template.ParseFS(tFS, filePath)
	if err != nil {
		return err
	}

	err = t.Execute(file, data)
	if err != nil {
		return err
	}

	return nil
}
