package pwa

import (
	"embed"
	"picow-led/internal/web"

	"github.com/labstack/echo/v4"
)

//go:embed templates
var templates embed.FS

type TemplateData struct {
	ServerPathPrefix string
	Version          string
	// AssetsToCache will auto prefixed with the `ServerPathPrefix`
	AssetsToCache []string
}

func Serve(e *echo.Echo, data TemplateData) {
	tFS := web.StripSubFromFS(templates)
	web.ServeGET(e, map[web.File]map[web.Key][]web.Value{
		"manifest.json": {
			"Content-Type": {
				"application/json",
			},
		},
		"service-worker.js": {
			"Content-Type": {
				"application/javascript",
			},
		},
	}, web.ServeData{
		FS:               tFS,
		ServerPathPrefix: data.ServerPathPrefix,
		TemplateData:     data,
	})
}

func Generate(path string, td TemplateData) error {
	tFS := web.StripSubFromFS(templates)

	err := web.GenerateFile(tFS, path, "manifest.json", td)
	if err != nil {
		return err
	}

	return web.GenerateFile(tFS, path, "service-worker.js", td)
}
