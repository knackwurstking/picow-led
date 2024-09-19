package endpoints

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func createApiEndpoints(e *echo.Echo) {
	// TODO: Create "api/colors", "api/devices" and "api/device" endpoints
	g := e.Group("/api")

	g.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, api)
	})

	createApiColorsEndpoints(g)
	createApiDevicesEndpoints(g)
	createApiDeviceEndpoints(g)
}
