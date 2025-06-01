package routes

import (
	"html/template"
	"io/fs"

	"github.com/labstack/echo/v4"
)

var (
	componentTemplates = []string{
		"components/device-list-item.go.html",
		"components/online-indicator.go.html",
		"components/svg-power.go.html",
	}
)

func ServeHTML(c echo.Context, f fs.FS, data any, patterns ...string) error {
	t, err := template.ParseFS(f, append(patterns, componentTemplates...)...)
	if err != nil {
		return err
	}

	resp := c.Response()
	resp.Header().Add("Content-Type", "text/html; charset=utf-8")

	return t.Execute(resp, data)
}
