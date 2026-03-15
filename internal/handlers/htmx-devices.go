package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"

	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/internal/htmx"
	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/knackwurstking/picow-led/internal/views/dialogs"
	"github.com/knackwurstking/picow-led/pkg/models"
)

func HTMXDevices(r *services.Registry) echo.HandlerFunc {
	log := env.NewLogger("handlers.HTMXDevices")

	return func(c echo.Context) error {
		devices, err := r.Device.List()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to retrieve devices: %v", err))
		}

		wg := sync.WaitGroup{}
		for _, d := range devices {
			wg.Go(func() {
				if color, err := r.Device.GetCurrentColor(d.ID); err != nil {
					// TODO: Pass errors to the frontend
					log.Warn("failed to get current color for device %d: %v", d.ID, err)
				} else {
					d.Color = color
				}
			})
		}
		wg.Wait()

		t := htmx.Devices(htmx.DevicesProps{
			Data: devices,
		})
		if err := t.Render(c.Request().Context(), c.Response()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to render template: %v", err))
		}
		return nil
	}
}

func HTMXToggleDevicePower(r *services.Registry) echo.HandlerFunc {
	log := env.NewLogger("handlers.HTMXToggleDevicePower")

	return func(c echo.Context) error {
		powerState, _ := strconv.ParseBool(strings.TrimSpace(c.FormValue("power_state")))

		deviceID, err := parseQueryID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid device ID: %v", err))
		}

		device, err := r.Device.Get(deviceID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to retrieve device: %v", err))
		}

		var color []uint8
		if c, err := r.Device.GetColor(device.ID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get current color: %v", err))
		} else {
			if !powerState {
				for range c {
					color = append(color, 0)
				}
			} else {
				color = c
			}
		}

		log.Debug("Toggling power state for device with %s to %s", device.Name, color)

		if err = r.Device.SetCurrentColor(device.ID, color); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to set current color: %v", err))
		}

		return nil
	}
}

func HTMXAddDeviceDialog(r *services.Registry, method string) echo.HandlerFunc {
	renderDialog := func(c echo.Context, open bool, formData dialogs.AddDeviceFormData, errs ...error) error {
		c.Response().Header().Set("HX-Trigger", "reload-devices")

		t := dialogs.AddDevice(dialogs.AddDeviceProps{
			AddDeviceFormData: formData,
			Open:              open,
			OOB:               true,
			Errors:            errs,
		})
		if err := t.Render(c.Request().Context(), c.Response()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to render template: %v", err))
		}
		return nil
	}

	switch method {
	case http.MethodGet:
		return func(c echo.Context) error {
			return renderDialog(c, true, dialogs.AddDeviceFormData{})
		}
	case http.MethodPost:
		return func(c echo.Context) error {
			var errs []error

			// Parse form values
			name := strings.TrimSpace(c.FormValue("name"))
			addr := strings.TrimSpace(c.FormValue("addr"))

			// Add to database
			device := models.NewDevice(addr, name)
			if _, err := r.Device.Add(device); err != nil {
				errs = append(errs, fmt.Errorf("failed to add device: %v", err))
			}

			open := len(errs) > 0
			formData := dialogs.AddDeviceFormData{
				Name: name,
				Addr: addr,
			}
			return renderDialog(c, open, formData, errs...)
		}
	}

	return nil
}

// TODO: Color handling, form value "color"
func HTMXEditDeviceDialog(r *services.Registry, method string) echo.HandlerFunc {
	parseForm := func(c echo.Context) (dialogs.EditDeviceFormData, []error) {
		var errs []error
		var formData dialogs.EditDeviceFormData

		id, err := parseQueryID(c)
		if err != nil {
			errs = append(errs, err)
		}

		formData.Addr = strings.TrimSpace(c.FormValue("addr"))
		formData.Name = strings.TrimSpace(c.FormValue("name"))
		formData.ID = id
		formData.Color = strings.TrimSpace(c.FormValue("color"))

		return formData, errs
	}

	renderDialog := func(c echo.Context, open bool, formData dialogs.EditDeviceFormData, errs ...error) error {
		c.Response().Header().Set("HX-Trigger", "reload-devices")

		t := dialogs.EditDevice(dialogs.EditDeviceProps{
			EditDeviceFormData: formData,
			Open:               open,
			OOB:                true,
			Errors:             errs,
		})
		if err := t.Render(c.Request().Context(), c.Response()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to render template: %v", err))
		}
		return nil
	}

	switch method {
	case http.MethodGet:
		return func(c echo.Context) error {
			var errs []error

			id, err := parseQueryID(c)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid device ID: %v", err))
			}

			// Get device from database
			device, err := r.Device.Get(id)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to retrieve device: %v", err))
			}

			formData := dialogs.EditDeviceFormData{
				ID:   id,
				Name: device.Name,
				Addr: device.Addr,
			}
			return renderDialog(c, true, formData, errs...)
		}

	case http.MethodPost:
		return func(c echo.Context) error {
			formData, errs := parseForm(c)

			if len(errs) == 0 {
				device := models.NewDevice(formData.Addr, formData.Name)
				device.ID = formData.ID

				if err := r.Device.Update(device); err != nil {
					errs = append(errs, fmt.Errorf("failed to update device: %v", err))
				} else {
					color := models.NewColorFromHex("", formData.Color)
					if err = r.Device.UpdateColor(device.ID, color.Duty...); err != nil {
						errs = append(errs, fmt.Errorf("failed to update device color: %v", err))
					}
				}
			}

			open := len(errs) > 0
			return renderDialog(c, open, formData, errs...)
		}

	case http.MethodDelete:
		return func(c echo.Context) error {
			formData, errs := parseForm(c)

			if len(errs) == 0 {
				if err := r.Device.Delete(formData.ID); err != nil {
					errs = append(errs, fmt.Errorf("failed to delete device: %v", err))
				}
			}

			open := len(errs) > 0
			return renderDialog(c, open, formData, errs...)
		}
	}

	return nil
}
