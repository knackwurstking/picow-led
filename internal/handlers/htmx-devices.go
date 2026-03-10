package handlers

import (
	"net/http"

	"github.com/knackwurstking/picow-led/internal/htmx"
	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/labstack/echo/v4"
)

func HTMXDevices(r *services.Registry) echo.HandlerFunc {
	return func(c echo.Context) error {
		t := htmx.Devices()
		if err := t.Render(c.Request().Context(), c.Response()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to render template: %v", err)
		}
		return nil
	}
}
