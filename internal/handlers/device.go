package handlers

import (
	"fmt"
	"net/http"

	"github.com/knackwurstking/picow-led/internal/pages"
	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/labstack/echo/v4"
)

func Device(r *services.Registry, method string) echo.HandlerFunc {
	render := func(c echo.Context) error {
		t := pages.Device()
		if err := t.Render(c.Request().Context(), c.Response()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Failed to render page: %w", err))
		}
		return nil
	}

	switch method {
	case http.MethodGet:
		return func(c echo.Context) error {
			return render(c)
		}
	}

	return nil
}
