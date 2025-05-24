package routes

import (
	"html/template"

	"github.com/labstack/echo/v4"
)

func RegisterPWA(e *echo.Echo, o *Options) {
	e.GET(o.ServerPathPrefix+"/manifest.json", func(c echo.Context) error {
		t, err := template.ParseFS(o.Templates, "pwa/manifest.json")
		if err != nil {
			return err
		}

		resp := c.Response()
		resp.Header().Add("Content-Type", "application/json; charset=utf-8")

		return t.Execute(resp, o.Global)
	})
}
