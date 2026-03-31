package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/a-h/templ"
	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/knackwurstking/picow-led/internal/templates/components/alert"
	"github.com/knackwurstking/picow-led/internal/templates/device"
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

		log.Debug("parsed pins from form: %#v", pins)
		return
	}

	render := func(c echo.Context, deviceID models.ID, deviceType models.DeviceType, pins ...uint8) error {
		t := device.PinsForm(device.PinsFormProps{
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
				fmt.Errorf("failed to render device pins: %w", err))
		}

		return nil
	}

	switch method {
	case http.MethodGet:
		return func(c echo.Context) error {
			id, err := utils.ParseQueryID(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest,
					fmt.Errorf("invalid device ID: %w", err))
			}

			device, err := r.Device.Get(id)
			if err != nil {
				return alert.RenderError(c, fmt.Sprintf("Failed to get device: %v", err))
			}

			pins, err := r.Device.GetPins(device.ID)
			if err != nil {
				return alert.RenderError(c, fmt.Sprintf("Failed to get device pins: %v", err))
			}

			return render(c, device.ID, device.Type, pins...)
		}
	case http.MethodPost:
		return func(c echo.Context) error {
			id, err := utils.ParseQueryID(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest,
					fmt.Errorf("invalid device ID: %w", err))
			}

			device, err := r.Device.Get(id)
			if err != nil {
				return alert.RenderError(c, fmt.Sprintf("Failed to get device: %v", err))
			}

			// Parse form
			pins, err := parse(c)
			if err != nil {
				return alert.RenderError(c, fmt.Sprintf("Invalid form data: %v", err))
			}

			// Update device pins
			if err = r.Device.SetPins(device.ID, pins); err != nil {
				return alert.RenderError(c, fmt.Sprintf("Failed to set device pins: %v", err))
			}

			return render(c, device.ID, device.Type, pins...)
		}
	}

	return nil
}
