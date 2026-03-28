package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/a-h/templ"
	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/knackwurstking/picow-led/internal/templates/components"
	"github.com/knackwurstking/picow-led/internal/utils"
	"github.com/knackwurstking/picow-led/pkg/models"
	"github.com/labstack/echo/v4"
)

func HTMXDevicePins(r *services.Registry, method string) echo.HandlerFunc {
	log := env.NewLogger("handlers.HTMXDevicePins")

	parse := func(c echo.Context) (pins []uint8, err error) {
		var formParams url.Values
		formParams, err = c.FormParams()
		if err != nil {
			return
		}

		// NOTE:`url.Values{"device-pins-pin":[]string{"0", "1", "2", "3"}, "id":[]string{"2"}}`
		if pinsValue, ok := formParams["device-pins-pin"]; ok {
			for _, pinStr := range pinsValue {
				var p uint64
				p, err = strconv.ParseUint(pinStr, 10, 8)
				if err != nil {
					return
				}

				pins = append(pins, uint8(p))
			}
		}

		log.Debug("Parsed pins from form: %#v", pins)
		return
	}

	render := func(c echo.Context, deviceID models.ID, deviceType models.DeviceType, pins ...uint8) error {
		t := components.DevicePins(components.DevicePinsProps{
			ID:         env.IDDevicePins,
			DeviceID:   deviceID,
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

			return render(c, device.ID, device.Type, pins...)
		}
	case http.MethodPost:
		return func(c echo.Context) error {
			_, err := parse(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest,
					fmt.Errorf("Invalid form data: %w", err))
			}

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

			// TODO: Parse form

			// TODO: Update device pins

			return render(c, device.ID, device.Type)
		}
	}

	return nil
}
