package handlers

import (
	"net/http"

	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/labstack/echo/v4"
)

func HTMXDevicePins(r *services.Registry, method string) echo.HandlerFunc {
	render := func(c echo.Context, pins ...uint8) error {
		// TODO: Render template with pins

		return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented yet")
	}

	switch method {
	case http.MethodGet:
		return func(c echo.Context) error {
			pins := []uint8{}

			// TODO: Get pins from the device registry

			return render(c, pins...)
		}
	}

	return nil
}
