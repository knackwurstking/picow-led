package routes

import (
	"net/http"
	"picow-led/internal/api"

	"github.com/labstack/echo/v4"
)

// apiDevices - GET - "/api/devices"
func apiDevices(e *echo.Echo, data Options) {
	e.GET(data.ServerPathPrefix+"/api/devices", func(c echo.Context) error {
		return c.JSON(http.StatusOK, api.GetDevices(data.Api))
	})
}
