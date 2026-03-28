package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"

	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/knackwurstking/picow-led/internal/templates/components"
	"github.com/knackwurstking/picow-led/internal/templates/components/dialogs"
	"github.com/knackwurstking/picow-led/internal/utils"
	"github.com/knackwurstking/picow-led/pkg/models"
)

func HTMXDevices(r *services.Registry) echo.HandlerFunc {
	log := env.NewLogger("handlers.HTMXDevices")

	return func(c echo.Context) error {
		var errs []error

		devices, err := r.Device.List()
		if err != nil {
			message := fmt.Sprintf("failed to retrieve devices: %v", err)
			log.Warn(message)
			errs = append(errs, errors.New(message))
		} else {
			wg := sync.WaitGroup{}
			for _, d := range devices {
				wg.Go(func() {
					if duty, err := r.Device.GetCurrentDuty(d.ID); err != nil {
						message := fmt.Sprintf("failed to retrieve current duty for device %d: %v", d.ID, err)
						log.Warn(message)
						errs = append(errs, errors.New(message))
					} else {
						d.Duty = duty
					}
				})
			}
			wg.Wait()
		}

		t := components.Devices(components.DevicesProps{
			Data:   devices,
			Errors: errs,
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
		var errs []error

		powerState, _ := strconv.ParseBool(strings.TrimSpace(c.FormValue("power_state")))

		deviceID, err := utils.ParseQueryID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%v: %s", err, c.QueryParam("id")))
		}

		if device, err := r.Device.Get(deviceID); err != nil {
			errs = append(errs, fmt.Errorf("failed to retrieve device: %w", err))
		} else {
			var duty []uint8
			if !powerState {
				for range device.Duty {
					duty = append(duty, 0)
				}
			} else {
				duty = device.Duty
			}

			log.Debug("Toggling power state for device with %s to %#v", device.Name, duty)

			if err = r.Device.SetCurrentDuty(device.ID, duty); err != nil {
				errs = append(errs, fmt.Errorf("failed to toggle power for device %d: %w", device.ID, err))
			}
		}

		// Handle errors (e.g. Render error messages)
		for _, err := range errs {
			t := components.OOBAddAlert(env.IDAlertContainer, components.AlertTypeError, err.Error())
			if err = t.Render(c.Request().Context(), c.Response()); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to render error alert: %w", err))
			}
		}

		return nil
	}
}

func HTMXAddDeviceDialog(r *services.Registry, method string) echo.HandlerFunc {
	log := env.NewLogger("handlers.HTMXAddDeviceDialog")

	parseForm := func(c echo.Context) (dialogs.AddDeviceFormData, []error) {
		var errs []error
		var formData dialogs.AddDeviceFormData

		formData.Addr = strings.TrimSpace(c.FormValue("addr"))
		formData.Name = strings.TrimSpace(c.FormValue("name"))
		formData.DeviceType = models.DeviceType(strings.TrimSpace(c.FormValue("device_type")))

		return formData, errs
	}

	renderDialog := func(c echo.Context, open bool, formData dialogs.AddDeviceFormData, errs ...error) error {
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

			formData, errs := parseForm(c)

			log.Debug("Adding device: %#v", formData)

			// Add device to database and set its color
			device := models.NewDevice(formData.Addr, formData.Name, formData.DeviceType)

			if id, err := r.Device.Add(device); err != nil {
				errs = append(errs, fmt.Errorf("failed to add device: %v", err))
			} else {
				device.ID = id
			}

			c.Response().Header().Set("HX-Trigger", "reload-devices")

			open := len(errs) > 0
			return renderDialog(c, open, formData, errs...)
		}
	}

	return nil
}

func HTMXEditDeviceDialog(r *services.Registry, method string) echo.HandlerFunc {
	log := env.NewLogger("handlers.HTMXEditDeviceDialog")

	parseForm := func(c echo.Context) (formData dialogs.EditDeviceFormData, errs []error) {
		id, err := utils.ParseQueryID(c)
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid device ID: %w", err))
		}
		formData.ID = id

		formData.Addr = strings.TrimSpace(c.FormValue("addr"))
		formData.Name = strings.TrimSpace(c.FormValue("name"))
		formData.DeviceType = models.DeviceType(strings.TrimSpace(c.FormValue("device_type")))

		return
	}

	renderDialog := func(c echo.Context, open bool, formData dialogs.EditDeviceFormData, errs ...error) error {
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

			id, err := utils.ParseQueryID(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%v: %s", err, c.QueryParam("id")))
			}

			formData := dialogs.EditDeviceFormData{
				ID:         id,
				Name:       "",
				Addr:       "",
				DeviceType: "",
			}

			// Get device from database
			if device, err := r.Device.Get(id); err != nil {
				errs = append(errs, fmt.Errorf("failed to retrieve device: %v", err))
			} else {
				log.Debug("Retrieved device: %#v", device)

				color := device.ToColor()
				log.Debug("Editing device with current color: %#v", color)

				formData.Name = device.Name
				formData.Addr = device.Addr
				formData.DeviceType = device.Type
			}

			return renderDialog(c, true, formData, errs...)
		}

	case http.MethodPost:
		return func(c echo.Context) error {
			formData, errs := parseForm(c)
			if len(errs) > 0 {
				return renderDialog(c, true, formData, errs...)
			}

			if device, err := r.Device.Get(formData.ID); err != nil {
				errs = append(errs, fmt.Errorf("failed to retrieve device for update: %v", err))
			} else {
				device.Addr = formData.Addr
				device.Name = formData.Name
				device.Type = formData.DeviceType

				log.Debug("Updating device with new data: %#v", device)

				if err = r.Device.Update(device); err != nil {
					errs = append(errs, fmt.Errorf("failed to update device: %v", err))
				}
			}

			c.Response().Header().Set("HX-Trigger", "reload-devices")

			return renderDialog(c, len(errs) > 0, formData, errs...)
		}

	case http.MethodDelete:
		return func(c echo.Context) error {
			id, err := utils.ParseQueryID(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%v: %s", err, c.QueryParam("id")))
			}

			if err := r.Device.Delete(id); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to delete device: %v", err))
			}

			c.Response().Header().Set("HX-Trigger", "reload-devices")

			return renderDialog(c, false, dialogs.EditDeviceFormData{}) // Close dialog
		}
	}

	return nil
}
