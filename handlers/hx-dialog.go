package handlers

import (
	"fmt"
	"net/http"

	"github.com/knackwurstking/picow-led/components"
	"github.com/knackwurstking/picow-led/models"
	"github.com/knackwurstking/picow-led/services"
	"github.com/labstack/echo/v4"
)

type HXDialogs struct {
	registry *services.Registry
}

func NewHXDialogs(r *services.Registry) *HXDialogs {
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

	var device *models.Device
	if deviceID > 0 {
		device, err = h.registry.Devices.Get(deviceID)
		if err != nil {
			if services.ErrNotFound == err {
				return echo.NewHTTPError(
					http.StatusBadRequest,
					NewValidationError("device with ID %d not found", deviceID),
				)
			}
			return NewDatabaseError("failed to fetch device with ID %d", deviceID)
		}
	}

	if device != nil {
		err = components.DialogEditDevice(device).Render(c.Request().Context(), c.Response())
		if err != nil {
			return err
		}
	} else {
		err = components.DialogNewDevice().Render(c.Request().Context(), c.Response())
		if err != nil {
			return err
		}
	}

	return err
}

// TODO: Find a way how to re render dialogs with error message, maybe do an oob render whatever
func (h HXDialogs) PostEditDevice(c echo.Context) error {
	device := h.parseEditDeviceForm(c)

	if !device.Validate() {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			NewValidationError("device validation failed, invalid form data %#v", device),
		)
	}

	if _, err := h.registry.Devices.Add(device); err != nil {
		return NewDatabaseError("failed to add device %s", device.Name)
	}

	c.Response().Header().Set("HX-Trigger", "reload")
	return nil
}

func (h HXDialogs) PutEditDevice(c echo.Context) error {
	// TODO: ...

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
	address := c.FormValue("address")
	port := c.FormValue("port")
	name := c.FormValue("name")

	return models.NewDevice(models.Addr(fmt.Sprintf("%s:%s", address, port)), name)
}
