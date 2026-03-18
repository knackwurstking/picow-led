package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"

	"github.com/knackwurstking/picow-led/internal/components"
	"github.com/knackwurstking/picow-led/internal/components/dialogs"
	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/knackwurstking/picow-led/internal/utils"
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

		t := components.Devices(components.DevicesProps{
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

		deviceID, err := utils.ParseQueryID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%v: %s", err, c.QueryParam("id")))
		}

		device, err := r.Device.Get(deviceID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to retrieve device: %v", err))
		}

		var color []uint8
		if !powerState {
			for range device.Color {
				color = append(color, 0)
			}
		} else {
			color = device.Color
		}

		log.Debug("Toggling power state for device with %s to %#v", device.Name, color)

		if err = r.Device.SetCurrentColor(device.ID, color); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to set current color: %v", err))
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

			formData, errs := parseForm(c)

			log.Debug("Adding device: %#v", formData)

			// Add device to database and set its color
			device := models.NewDevice(formData.Addr, formData.Name, formData.DeviceType)

			id, err := r.Device.Add(device)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to add device: %v", err))
			}
			device.ID = id

			open := len(errs) > 0
			return renderDialog(c, open, formData, errs...)
		}
	}

	return nil
}

func HTMXEditDeviceDialog(r *services.Registry, method string) echo.HandlerFunc {
	log := env.NewLogger("handlers.HTMXEditDeviceDialog")

	parseForm := func(c echo.Context) (dialogs.EditDeviceFormData, []error) {
		var errs []error
		var formData dialogs.EditDeviceFormData

		id, err := utils.ParseQueryID(c)
		if err != nil {
			errs = append(errs, fmt.Errorf("%v: %s", err, c.QueryParam("id")))
		}
		formData.ID = id

		formData.Addr = strings.TrimSpace(c.FormValue("addr"))
		formData.Name = strings.TrimSpace(c.FormValue("name"))
		formData.DeviceType = models.DeviceType(strings.TrimSpace(c.FormValue("device_type")))

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

			id, err := utils.ParseQueryID(c)
			if err != nil {
				errs = append(errs, fmt.Errorf("%v: %s", err, c.QueryParam("id")))
			}

			// Get device from database
			device, err := r.Device.Get(id)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to retrieve device: %v", err))
			}
			log.Debug("Retrieved device: %#v", device)

			color := device.ToColor()
			log.Debug("Editing device with current color: %#v", color)

			formData := dialogs.EditDeviceFormData{
				ID:         id,
				Name:       device.Name,
				Addr:       device.Addr,
				DeviceType: device.Type,
			}

			return renderDialog(c, true, formData, errs...)
		}

	case http.MethodPost:
		return func(c echo.Context) error {
			formData, errs := parseForm(c)

			if len(errs) == 0 {
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
			}

			open := len(errs) > 0
			return renderDialog(c, open, formData, errs...)
		}

	case http.MethodDelete:
		return func(c echo.Context) error {
			formData, errs := parseForm(c)

			log.Debug("Deleting device with ID %d", formData.ID)

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
