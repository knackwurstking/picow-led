package routes

import (
	"log/slog"

	"github.com/labstack/echo/v4"
)

func RegisterFrontend(e *echo.Echo, o *Options) {
	devicesHandlerFunc := func(c echo.Context) error {
		err := ServeHTML(c, o.Templates, o.Devices(),
			"main.go.html",
			"layouts/base.go.html",
			"content/devices.go.html",
		)
		if err != nil {
			slog.Error(err.Error(), "path", c.Request().URL.Path)
		}
		return err
	}

	e.GET(o.ServerPathPrefix+"/", devicesHandlerFunc)
	e.GET(o.ServerPathPrefix+"/devices", devicesHandlerFunc)

	// TODO: ...
}
