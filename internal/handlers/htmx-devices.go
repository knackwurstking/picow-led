package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/knackwurstking/picow-led/internal/htmx"
	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/knackwurstking/picow-led/internal/views/dialogs"
	"github.com/knackwurstking/picow-led/pkg/models"
)

func HTMXDevices(r *services.Registry) echo.HandlerFunc {
	return func(c echo.Context) error {
		devices, err := r.Device.List()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to retrieve devices: %v", err))
		}

		t := htmx.Devices(htmx.DevicesProps{
			Data: devices,
		})
		if err := t.Render(c.Request().Context(), c.Response()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to render template: %v", err))
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
