package devices

import (
	"log/slog"
	"net/http"

	"github.com/knackwurstking/picow-led/handlers/components/oob"
	"github.com/knackwurstking/picow-led/handlers/home/components"
	"github.com/knackwurstking/picow-led/services"
	"github.com/knackwurstking/picow-led/utils"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	registry *services.Registry
}

func NewHandler(r *services.Registry) *Handler {
	return &Handler{
		registry: r,
	}
}

func (h *Handler) Register(e *echo.Echo) {
	utils.Register(e, http.MethodGet,
		"/htmx/home/devices", h.GetDevices)
	utils.Register(e, http.MethodDelete,
		"/htmx/home/devices/delete", h.DeleteDevice)
	utils.Register(e, http.MethodPost,
		"/htmx/home/devices/toggle-power", h.PostTogglePowerDevice)
}

func (h *Handler) GetDevices(c echo.Context) error {
	slog.Info("Render devices section for the home page")

	// Get devices...
	devices, err := h.registry.Devices.List()
	if err != nil {
		return utils.WrapError(err, "list devices")
	}

	rDevices, err := services.ResolveDevices(h.registry, devices...)
	if err != nil {
		return utils.WrapError(err, "resolve devices")
	}

	return components.SectionDevices(false, rDevices).Render(c.Request().Context(), c.Response())
}

func (h *Handler) DeleteDevice(c echo.Context) error {
	deviceID, err := utils.QueryParamDeviceID(c, "id", false)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.WrapError(err, "get device ID from query parameter"))
	}

	slog.Info("Delete a device", "id", deviceID)

	if err = h.registry.Devices.Delete(deviceID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.WrapError(err, "delete device"))
	}

	c.Response().Header().Set("HX-Trigger", "reloadDevices")
	return nil
}

func (h *Handler) PostTogglePowerDevice(c echo.Context) error {
	deviceID, err := utils.QueryParamDeviceID(c, "id", false)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.WrapError(err, "get device ID from query parameter"))
	}

	slog.Info("Toggle power for a device", "id", deviceID)

	color, err := h.registry.DeviceControls.TogglePower(deviceID)
	if err != nil {
		err = utils.WrapError(err, "toggle power for device %d", deviceID)
		oob.OOBRenderPageHomeDeviceError(c, deviceID, err)
		oob.OOBRenderPageHomeDevicePowerButton(c, deviceID, color)

		if services.IsNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	oob.OOBRenderPageHomeDeviceError(c, deviceID, nil)
	oob.OOBRenderPageHomeDevicePowerButton(c, deviceID, color)

	return nil
}
