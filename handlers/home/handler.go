package home

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/knackwurstking/picow-led/handlers/home/components"
	"github.com/knackwurstking/picow-led/handlers/utils"
	"github.com/knackwurstking/picow-led/services"
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

// TODO: Either move device and group methods to separate handlers or rename routes
func (h *Handler) Register(e *echo.Echo) {
	utils.Register(e, http.MethodGet, "", h.GetHomePage)

	utils.Register(e, http.MethodGet, "/htmx/home/section/devices", h.GetSectionDevices)
	utils.Register(e, http.MethodGet, "/htmx/home/section/groups", h.GetSectionGroups)

	utils.Register(e, http.MethodDelete, "/htmx/devices/delete", h.DeleteDevice)
	utils.Register(e, http.MethodPost, "/htmx/devices/toggle-power", h.PostDeviceTogglePower)

	utils.Register(e, http.MethodDelete, "/htmx/groups/delete", h.DeleteGroup)
	utils.Register(e, http.MethodPost, "/htmx/groups/toggle-power", h.PostGroupTogglePower)
}

func (h *Handler) GetHomePage(c echo.Context) error {
	err := components.PageHome().Render(c.Request().Context(), c.Response())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return nil
}

func (h *Handler) GetSectionDevices(c echo.Context) error {
	slog.Info("Render devices section for the home page")

	// Get devices...
	devices, err := h.registry.Devices.List()
	if err != nil {
		return fmt.Errorf("failed to list devices: %v", err)
	}

	rDevices, err := services.ResolveDevices(h.registry, devices...)
	if err != nil {
		return fmt.Errorf("failed to resolve devices: %v", err)
	}

	return components.SectionDevices(false, rDevices).Render(c.Request().Context(), c.Response())
}

func (h *Handler) GetSectionGroups(c echo.Context) error {
	slog.Info("Render groups section for the home page")

	// Get groups...
	groups, err := h.registry.Groups.List()
	if err != nil {
		return fmt.Errorf("failed to list groups: %v", err)
	}

	// ...resolve them
	resolvedGroups, err := services.ResolveGroups(h.registry, groups...)
	if err != nil {
		return fmt.Errorf("failed to resolve groups: %v", err)
	}

	return components.SectionGroups(false, resolvedGroups).Render(c.Request().Context(), c.Response())
}

func (h *Handler) DeleteDevice(c echo.Context) error {
	deviceID, err := utils.QueryParamDeviceID(c, "id", false)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	slog.Info("Delete a device", "id", deviceID)

	if err = h.registry.Devices.Delete(deviceID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.Response().Header().Set("HX-Trigger", "reloadDevices")
	return nil
}

// TODO: Delete a group method

func (h *Handler) PostDeviceTogglePower(c echo.Context) error {
	deviceID, err := utils.QueryParamDeviceID(c, "id", false)
	if err != nil {
		return fmt.Errorf("Failed to get device id from query parameter: %s", err.Error())
	}

	slog.Info("Toggle power for a device", "id", deviceID)

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

// TODO: Post device toggle power method
