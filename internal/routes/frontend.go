package routes

import (
	"picow-led/internal/api"
	"picow-led/ui"

	"github.com/labstack/echo/v4"
)

var FrontendCache []*api.Device

func frontend(e *echo.Echo, data Options) {
	e.GET(data.ServerPathPrefix+"", func(c echo.Context) error {
		return ui.DevicesPage(data.ServerPathPrefix, FrontendCache...).Render(c.Response().Writer)
	})

	e.GET(data.ServerPathPrefix+"/settings", func(c echo.Context) error {
		return ui.SettingsPage(data.ServerPathPrefix, FrontendCache...).Render(c.Response().Writer)
	})
}
