package routes

import (
	"fmt"
	"net/http"
	"net/url"
	"picow-led/internal/api"
	"picow-led/ui"

	"github.com/labstack/echo/v4"
)

var FrontendCache []*api.Device

func frontend(e *echo.Echo, data Options) {
	e.GET(data.ServerPathPrefix+"", func(c echo.Context) error {
		return ui.DevicesPage(data.ServerPathPrefix, FrontendCache...).Render(c.Response().Writer)
	})

	e.GET(data.ServerPathPrefix+"/devices/:addr", func(c echo.Context) error {
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

		return ui.DevicesAddrPage(data.ServerPathPrefix, device).Render(c.Response().Writer)
	})

	e.GET(data.ServerPathPrefix+"/settings", func(c echo.Context) error {
		return ui.SettingsPage(data.ServerPathPrefix).Render(c.Response().Writer)
	})
}
