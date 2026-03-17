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
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Invalid device ID: %w", err))
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
