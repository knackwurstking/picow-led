package routes

import (
	"html/template"

	"github.com/labstack/echo/v4"
)

const manifest = `
{
    "version": "v1.0.0",
    "manifest_version": 3,
    "name": "PicoW LED",
    "short_name": "picow-led",
    "id": "",
    "icons": [
        {
            "src": "/static/icons/pwa-64x64.png",
            "sizes": "64x64",
            "type": "image/png"
        },
        {
            "src": "/static/icons/pwa-192x192.png",
            "sizes": "192x192",
            "type": "image/png"
        },
        {
            "src": "/static/icons/pwa-512x512.png",
            "sizes": "512x512",
            "type": "image/png"
        },
        {
            "src": "/static/icons/maskable-icon-512x512.png",
            "sizes": "512x512",
            "type": "image/png",
            "purpose": "maskable"
        }
    ],
    "screenshots": [
        {
            "src": "/static/screenshots/626x338.png",
            "sizes": "626x338",
            "type": "image/png",
            "form_factor": "wide",
            "label": "App Preview"
        },
        {
            "src": "/static/screenshots/328x626.png",
            "sizes": "328x626",
            "type": "image/png",
            "form_factor": "narrow",
            "label": "App Preview"
        }
    ],
    "theme_color": "#09090b",
    "background_color": "#09090b",
    "display": "standalone",
    "scope": ".",
    "start_url": "./",
    "public_path": ""
}
`

func pwa(e *echo.Echo, data Options) {
	e.GET(data.ServerPathPrefix+"/manifest.json", func(c echo.Context) error {
		t, err := template.New("manifest.json").Parse(manifest)
		if err != nil {
			return err
		}
		c.Response().Header().Add("Content-Type", "application/json")
		return t.Execute(c.Response().Writer, data)
	})

	//e.GET(data.ServerPathPrefix+"/js/service-worker.js", func(c echo.Context) error {
	//})
}
