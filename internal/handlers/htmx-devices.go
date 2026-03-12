package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/knackwurstking/picow-led/internal/htmx"
	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/knackwurstking/picow-led/internal/views/dialogs"
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
	renderDialog := func(open bool, formData dialogs.AddDeviceFormData, errs ...error) error {
		return nil
	}

	switch method {
	case http.MethodGet:
		return func(c echo.Context) error {
			return renderDialog(true, dialogs.AddDeviceFormData{})
		}
	case http.MethodPost:
		return func(c echo.Context) error {
			// TODO: Parse form data and close the dialog, or render with errors..
			var errs []error
			name := strings.TrimSpace(c.FormValue("name"))
			// ...
			return renderDialog(false, dialogs.AddDeviceFormData{}, errs...)
		}
	}

	return nil
}
