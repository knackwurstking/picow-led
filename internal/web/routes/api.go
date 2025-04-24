package routes

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

// apiDevices - GET - "/api/devices"
func apiDevices(e *echo.Echo, data Options) {
	e.GET(data.ServerPathPrefix+"/api/devices", func(c echo.Context) error {
		// TODO: ...

		return fmt.Errorf("under construction")
	})
}
