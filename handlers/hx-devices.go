package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/knackwurstking/picow-led/services"
	"github.com/labstack/echo/v4"
)

type HxDevices struct {
	registry *services.Registry
}

func NewHxDevices(registry *services.Registry) *HxDevices {
	return &HxDevices{
		registry: registry,
	}
}

func (h *HxDevices) Register(e *echo.Echo) {
	Register(e, http.MethodDelete, "/htmx/devices/delete", h.Delete)
	Register(e, http.MethodPost, "/htmx/devices/toggle-power", h.PostTogglePower)
}

func (h *HxDevices) Delete(c echo.Context) error {
	deviceID, err := QueryParamDeviceID(c, "id", false)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err = h.registry.Devices.Delete(deviceID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.Response().Header().Set("HX-Trigger", "reload")
	return nil
}

func (h *HxDevices) PostTogglePower(c echo.Context) error {
	deviceID, err := QueryParamDeviceID(c, "id", false)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Errorf("Failed to toggle power for device %d: %s", deviceID, err.Error()),
		)
	}

	color, err := h.registry.DeviceControls.TogglePower(deviceID)
	if err != nil {
		err = fmt.Errorf("Failed to toggle power for device %d: %s", deviceID, err.Error())
		OOBRenderPageHomeDeviceError(c, deviceID, err)

		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	data, _ := json.Marshal(struct {
		Color []uint8 `json:"color"`
	}{
		Color: color,
	})

	// Remove any existing error message
	OOBRenderPageHomeDeviceError(c, deviceID, nil)

	c.Response().Header().Set("HX-Trigger-After-Settle", string(data))
	// TODO: Set the correct trigger, right now there is no trigger needed
	//c.Response().Header().Set("HX-Trigger", "reload")

	return nil
}
