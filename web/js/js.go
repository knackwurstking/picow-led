package js

import (
	"embed"
	"picow-led/internal/web"

	"github.com/labstack/echo/v4"
)

//go:embed templates
var templates embed.FS

type TemplateData struct {
	ServerPathPrefix string
}

func Serve(e *echo.Echo, data TemplateData) {
	tFS := web.StripSubFromFS(templates)
	web.ServeGET(e, map[web.File]map[web.Key][]web.Value{
		"main.js": {
			"Content-Type": {
				"application/javascript",
			},
		},
		"api.js": {
			"Content-Type": {
				"application/javascript",
			},
		},
		"page-devices.js": {
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

	return web.GenerateFile(tFS, path, "main.js", td)
}
