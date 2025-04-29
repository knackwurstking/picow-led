package routes

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"net/url"
	"picow-led/internal/api"

	"github.com/labstack/echo/v4"
)

const (
	contentDevices     content = "devices"
	contentDevicesAddr content = "devices-addr"
	contentSettings    content = "settings"
)

var FrontendCache []*api.Device

type content string

type Frontend struct {
	ServerPathPrefix string
	Templates        fs.FS
}

func (f *Frontend) BasicPatterns() []string {
	return []string{
		"components/online-indicator.go.html",
		"components/power-button.go.html",
		"scripts/base-layout.js",
		"scripts/window-api.js",
		"scripts/window-utils.js",
		"scripts/window-ws.js",
	}
}

// serve template data
func (f *Frontend) serve(c echo.Context, content content, data frontendTemplateData) error {
	patterns := []string{
		"page.go.html",               // There is only one page for now
		"layout/base-layout.go.html", // There is also only on layout for now
		fmt.Sprintf("layout/content/%s.go.html", content),
	}
	patterns = append(patterns, f.BasicPatterns()...)

	t, err := template.ParseFS(f.Templates, patterns...)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Add("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(c.Response().Writer, data)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return nil
}

type frontendTemplateData struct {
	ServerPathPrefix string
	Title            string
}

func frontend(e *echo.Echo, o Frontend) {
	e.GET(o.ServerPathPrefix+"/", func(c echo.Context) error {
		return o.serve(c, contentDevices, frontendTemplateData{
			ServerPathPrefix: o.ServerPathPrefix,
			Title:            "PicoW LED | Devices",
		})
	})

	e.GET(o.ServerPathPrefix+"/devices/:addr", func(c echo.Context) error {
		addr, err := url.QueryUnescape(c.Param("addr"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		var device *api.Device
		for _, d := range FrontendCache {
			if d.Server.Addr == addr {
				device = d
				break
			}
		}

		if device == nil {
			return c.String(http.StatusNotFound, fmt.Sprintf("device \"%s\" not found", addr))
		}

		return o.serve(c, contentDevicesAddr, frontendTemplateData{
			ServerPathPrefix: o.ServerPathPrefix,
			Title:            fmt.Sprintf("PicoW LED | %s", addr),
		})
	})

	e.GET(o.ServerPathPrefix+"/settings", func(c echo.Context) error {
		return o.serve(c, contentSettings, frontendTemplateData{
			ServerPathPrefix: o.ServerPathPrefix,
			Title:            "PicoW LED | Settings",
		})
	})
}
