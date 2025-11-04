package handlers

import (
	"fmt"
	"log/slog"
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
	slog.Info("Toggle power for device")

	deviceID, err := QueryParamDeviceID(c, "id", false)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Errorf("Failed to get device id from query parameter: %s", err.Error()),
		)
	}

	color, err := h.registry.DeviceControls.TogglePower(deviceID)
	if err != nil {
		err = fmt.Errorf("Failed to toggle power for device %d: %s", deviceID, err.Error())
		OOBRenderPageHomeDeviceError(c, deviceID, err)
		OOBRenderPageHomeDevicePowerButton(c, deviceID, color)

		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	OOBRenderPageHomeDeviceError(c, deviceID, nil)
	OOBRenderPageHomeDevicePowerButton(c, deviceID, color)

	return nil
}
