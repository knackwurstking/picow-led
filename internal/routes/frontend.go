package routes

import (
	"picow-led/internal/api"
	"picow-led/ui"

	"github.com/labstack/echo/v4"
)

func frontend(e *echo.Echo, data Options) {
	// Page Data: "/" - devices
	devices := []*api.Device{}
	for _, s := range data.Api.Servers {
		devices = append(devices, &api.Device{Server: s})
	}

	e.GET(data.ServerPathPrefix+"", func(c echo.Context) error {
		return ui.DevicesPage(data.ServerPathPrefix, devices...).Render(c.Response().Writer)
	})

	e.GET(data.ServerPathPrefix+"/settings", func(c echo.Context) error {
		// TODO: Settings page missing
		return ui.DevicesPage(data.ServerPathPrefix, devices...).Render(c.Response().Writer)
	})
}
