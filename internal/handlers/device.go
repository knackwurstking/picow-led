package handlers

import (
	"net/http"

	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/labstack/echo/v4"
)

func Device(r *services.Registry) echo.HandlerFunc {
	// TODO: ...
	return func(c echo.Context) error {
		// ...

		return echo.NewHTTPError(http.StatusNotImplemented, "Work in Progress")
	}
}
