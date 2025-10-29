package handlers

import (
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
	Register(e, http.MethodGet, "/htmx/devices/delete", h.Delete)
}

func (h *HxDevices) Delete(c echo.Context) error {
	deviceID, err := QueryParamDeviceID(c, "id", false)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err = h.registry.Devices.Delete(deviceID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return nil
}
