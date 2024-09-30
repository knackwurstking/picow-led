package endpoints

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/knackwurstking/picow-led-server/pkg/api"
)

func createApiEndpoints(e *echo.Echo, a *api.API, changeCallback func()) {
	g := e.Group("/api")

	g.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, a)
	})

	createApiColorsEndpoints(g, a, changeCallback)
	createApiDevicesEndpoints(g, a)
	createApiDeviceEndpoints(g, a, changeCallback)
}
