package endpoints

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/knackwurstking/picow-led-server/pkg/api"
	"github.com/knackwurstking/picow-led-server/pkg/clients"
)

func createApiEndpoints(e *echo.Echo, c *clients.Clients, a *api.API, apiChangeCallback func()) {
	g := e.Group("/api")

	g.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, a)
	})

	createApiColorsEndpoints(g, a, apiChangeCallback)
	createApiDevicesEndpoints(g, a)
	createApiDeviceEndpoints(g, c, a, apiChangeCallback)
}
