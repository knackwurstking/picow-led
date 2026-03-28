package handlers

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/knackwurstking/picow-led/internal/templates/components"
	"github.com/knackwurstking/picow-led/internal/utils"
	"github.com/knackwurstking/picow-led/pkg/models"
	"github.com/labstack/echo/v4"
)

func HTMXDevicePins(r *services.Registry, method string) echo.HandlerFunc {
	render := func(c echo.Context, deviceType models.DeviceType, pins ...uint8) error {
		t := components.DevicePins(components.DevicePinsProps{
			ID:         env.IDDevicePins,
			DeviceType: deviceType,
			Pins:       pins,
			Attributes: templ.Attributes{
				"hx-swap-oob": "true",
			},
		})

		if err := t.Render(c.Request().Context(), c.Response()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError,
				fmt.Errorf("Failed to render device pins: %w", err))
		}

		return nil
	}

	switch method {
	case http.MethodGet:
		return func(c echo.Context) error {
			id, err := utils.ParseQueryID(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest,
					fmt.Errorf("Invalid device ID: %w", err))
			}

			device, err := r.Device.Get(id)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError,
					fmt.Errorf("Failed to get device: %w", err))
			}

			pins, err := r.Device.GetPins(device.ID)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError,
					fmt.Errorf("Failed to get device pins: %w", err))
			}

			return render(c, device.Type, pins...)
		}
	}

	return nil
}
