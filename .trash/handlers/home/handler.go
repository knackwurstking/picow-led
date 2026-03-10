package home

import (
	"log/slog"
	"net/http"
	"slices"
	"sync"

	"github.com/knackwurstking/picow-led/components/oob"
	"github.com/knackwurstking/picow-led/handlers/home/templates"
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
	utils.Register(e, http.MethodGet, "", h.GetPage)
}

func (h *Handler) GetPage(c echo.Context) error {
	err := templates.PageHome().Render(c.Request().Context(), c.Response())
	if err != nil {
		return utils.WrapError(err, "render home page")
	}

	return nil
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

	return templates.SectionDevices(false, rDevices).Render(c.Request().Context(), c.Response())
}

func (h *Handler) GetGroups(c echo.Context) error {
	slog.Info("Render groups section for the home page")

	// Get groups...
	groups, err := h.registry.Groups.List()
	if err != nil {
		return utils.WrapError(err, "list groups")
	}

	// ...resolve them
	resolvedGroups, err := services.ResolveGroups(h.registry, groups...)
	if err != nil {
		return utils.WrapError(err, "resolve groups")
	}

	return templates.SectionGroups(false, resolvedGroups).Render(c.Request().Context(), c.Response())
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

func (h *Handler) DeleteGroup(c echo.Context) error {
	groupID, err := utils.QueryParamGroupID(c, "id", false)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.WrapError(err, "get group ID from query parameter"))
	}

	slog.Info("Delete a group", "id", groupID)

	if err = h.registry.Groups.Delete(groupID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.WrapError(err, "delete group"))
	}

	c.Response().Header().Set("HX-Trigger", "reloadGroups")
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

func (h *Handler) PostTurnOnGroup(c echo.Context) error {
	groupID, err := utils.QueryParamGroupID(c, "id", false)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.WrapError(err, "get group ID from query parameter"))
	}

	group, err := h.registry.Groups.Get(groupID)
	if err != nil {
		oob.OOBRenderPageHomeGroupError(c, groupID, []error{err})

		if services.IsNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, utils.WrapError(err, "group not found"))
		}

		return echo.NewHTTPError(http.StatusInternalServerError, utils.WrapError(err, "get group"))
	}

	slog.Info("Turn on a group", "id", groupID, "devices", group.Devices)

	devices, err := h.registry.Devices.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.WrapError(err, "list devices"))
	}

	wg := &sync.WaitGroup{}
	errs := make([]error, 0)
	for _, d := range devices {
		wg.Go(func() {
			if !slices.Contains(group.Devices, d.ID) {
				if err := h.registry.DeviceControls.TurnOff(d.ID); err != nil {
					if device, err2 := h.registry.Devices.Get(d.ID); err2 != nil {
						errs = append(errs, utils.WrapError(err2, "get device %d from the database", d.ID))
					} else {
						errs = append(errs, utils.WrapError(err, "turn off device \"%s\", which is not in this group", device.Name))
					}
				}

				return
			}

			if err := h.registry.DeviceControls.TurnOn(d.ID); err != nil {
				if device, err2 := h.registry.Devices.Get(d.ID); err2 != nil {
					errs = append(errs, utils.WrapError(err2, "get device %d from the database", d.ID))
				} else {
					errs = append(errs, utils.WrapError(err, "turn on device \"%s\"", device.Name))
				}
			}
		})
	}
	wg.Wait()

	oob.OOBRenderPageHomeGroupError(c, groupID, errs)

	return nil
}

func (h *Handler) PostTurnOffGroup(c echo.Context) error {
	groupID, err := utils.QueryParamGroupID(c, "id", false)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.WrapError(err, "get group ID from query parameter"))
	}

	group, err := h.registry.Groups.Get(groupID)
	if err != nil {
		oob.OOBRenderPageHomeGroupError(c, groupID, []error{err})

		if services.IsNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, utils.WrapError(err, "group not found"))
		}

		return echo.NewHTTPError(http.StatusInternalServerError, utils.WrapError(err, "get group"))
	}

	slog.Info("Turn off a group", "id", groupID, "devices", group.Devices)

	devices, err := h.registry.Devices.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.WrapError(err, "list devices"))
	}

	wg := &sync.WaitGroup{}
	errs := make([]error, 0)
	for _, d := range devices {
		wg.Go(func() {
			if err := h.registry.DeviceControls.TurnOff(d.ID); err != nil {
				if device, err2 := h.registry.Devices.Get(d.ID); err2 != nil {
					errs = append(errs, utils.WrapError(err2, "get device %d from the database", d.ID))
				} else {
					errs = append(errs, utils.WrapError(err, "turn off device \"%s\"", device.Name))
				}
			}
		})
	}
	wg.Wait()

	oob.OOBRenderPageHomeGroupError(c, groupID, errs)

	return nil
}
