package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/knackwurstking/picow-led/components"
	"github.com/knackwurstking/picow-led/models"
	"github.com/knackwurstking/picow-led/services"
	"github.com/labstack/echo/v4"
)

type HXDialogs struct {
	registry *services.Registry
}

func NewHxDialogs(r *services.Registry) *HXDialogs {
	return &HXDialogs{
		registry: r,
	}
}

func (h HXDialogs) Register(e *echo.Echo) {
	// Edit Device
	Register(e, http.MethodGet, "/htmx/dialog/edit-device", h.GetEditDevice)
	Register(e, http.MethodPost, "/htmx/dialog/edit-device", h.PostEditDevice)
	Register(e, http.MethodPut, "/htmx/dialog/edit-device", h.PutEditDevice)

	// Edit Group
	Register(e, http.MethodGet, "/htmx/dialog/edit-group", h.GetEditGroup)
	Register(e, http.MethodPost, "/htmx/dialog/edit-group", h.PostEditGroup)
	Register(e, http.MethodPut, "/htmx/dialog/edit-group", h.PutEditGroup)
}

func (h HXDialogs) GetEditDevice(c echo.Context) error {
	deviceID, err := QueryParamDeviceID(c, "id", true)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	slog.Info("Get dialog for edit or create a device", "id", deviceID)

	var device *models.Device
	if deviceID > 0 {
		device, err = h.registry.Devices.Get(deviceID)
		if err != nil {
			if services.IsNotFoundError(err) {
				return echo.NewHTTPError(http.StatusBadRequest,
					fmt.Errorf("device with ID %d not found", deviceID))
			}

			return fmt.Errorf("failed to fetch device with ID %d", deviceID)
		}
	}

	if device != nil {
		slog.Info("Device found, rendering edit dialog")
		if err = components.DialogEditDevice(device, false, nil).Render(
			c.Request().Context(), c.Response(),
		); err != nil {
			return fmt.Errorf("failed to render dialog: %v", err)
		}
	} else {
		slog.Info("Device not found, rendering new device dialog")
		if err = components.DialogNewDevice(false, nil).Render(
			c.Request().Context(), c.Response(),
		); err != nil {
			return fmt.Errorf("failed to render dialog: %v", err)
		}
	}

	return nil
}

func (h HXDialogs) PostEditDevice(c echo.Context) error {
	device := h.parseEditDeviceForm(c)

	if !device.Validate() {
		validationError := fmt.Errorf("device validation failed, invalid form data %#v", device)

		if err := components.DialogNewDevice(true, validationError).Render(
			c.Request().Context(), c.Response(),
		); err != nil {
			return fmt.Errorf("failed to render dialog: %v", err)
		}

		return echo.NewHTTPError(http.StatusBadRequest, validationError)
	}

	slog.Info("Submit handler from the create device dialog",
		"name", device.Name, "addr", device.Addr)

	if _, err := h.registry.Devices.Add(device); err != nil {
		databaseError := fmt.Errorf("failed to add device %s", device.Name)

		if err := components.DialogNewDevice(true, databaseError).Render(
			c.Request().Context(), c.Response(),
		); err != nil {
			return fmt.Errorf("failed to render dialog: %v", err)
		}

		return databaseError
	}

	c.Response().Header().Set("HX-Trigger", "reload")
	return nil
}

func (h HXDialogs) PutEditDevice(c echo.Context) error {
	deviceID, err := QueryParamDeviceID(c, "id", true)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	device := h.parseEditDeviceForm(c)
	device.ID = deviceID

	if !device.Validate() {
		validationError := fmt.Errorf("device validation failed, invalid form data %#v", device)

		if err := components.DialogEditDevice(device, true, validationError).Render(
			c.Request().Context(), c.Response(),
		); err != nil {
			return err
		}

		return echo.NewHTTPError(http.StatusBadRequest, validationError)
	}

	slog.Info("Submit handler from the edit device dialog",
		"id", device.ID, "name", device.Name, "addr", device.Addr)

	if err := h.registry.Devices.Update(device); err != nil {
		databaseError := fmt.Errorf("failed to update device %s", device.Name)

		if err := components.DialogEditDevice(device, true, databaseError).Render(
			c.Request().Context(), c.Response(),
		); err != nil {
			return fmt.Errorf("failed to render dialog: %v", err)
		}

		return databaseError
	}

	c.Response().Header().Set("HX-Trigger", "reload")
	return nil
}

func (h HXDialogs) GetEditGroup(c echo.Context) error {
	// TODO: ...

	return nil
}

func (h HXDialogs) PostEditGroup(c echo.Context) error {
	// TODO: ...

	c.Response().Header().Set("HX-Trigger", "reload")
	return nil
}
func (h HXDialogs) PutEditGroup(c echo.Context) error {
	// TODO: ...

	c.Response().Header().Set("HX-Trigger", "reload")
	return nil
}

func (h *HXDialogs) parseEditDeviceForm(c echo.Context) *models.Device {
	host := c.FormValue("host")
	port := c.FormValue("port")
	name := c.FormValue("device-name")

	return models.NewDevice(models.Addr(fmt.Sprintf("%s:%s", host, port)), name)
}
