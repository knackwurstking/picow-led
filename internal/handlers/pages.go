package handlers

import (
	"fmt"
	"net/http"

	"github.com/knackwurstking/picow-led/internal/pages"
	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/knackwurstking/picow-led/internal/utils"
	"github.com/knackwurstking/picow-led/pkg/models"
	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	t := pages.Home()
	if err := t.Render(c.Request().Context(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to render template: %v", err))
	}
	return nil
}

func Device(r *services.Registry, method string) echo.HandlerFunc {
	render := func(c echo.Context, device *models.Device) error {
		t := pages.Device(device)
		if err := t.Render(c.Request().Context(), c.Response()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Failed to render page: %w", err))
		}
		return nil
	}

	switch method {
	case http.MethodGet:
		return func(c echo.Context) error {
			id, err := utils.ParseQueryID(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%v: %s", err, c.QueryParam("id")))
			}

			device, err := r.Device.Get(id)
			if err != nil {
				return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Device not found: %w", err))
			}

			return render(c, device)
		}
	}

	return nil
}
