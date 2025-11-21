package home

import (
	"log/slog"
	"net/http"
	"slices"
	"sync"

	"github.com/knackwurstking/picow-led/errors"
	"github.com/knackwurstking/picow-led/handlers/components/oob"
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

func (h *Handler) Register(e *echo.Echo) {
	utils.Register(e, http.MethodGet, "", h.GetPage)
}

func (h *Handler) GetPage(c echo.Context) error {
	err := components.PageHome().Render(c.Request().Context(), c.Response())
	if err != nil {
		return errors.Wrap(err, "failed to render home page")
	}

	return nil
}

func (h *Handler) GetDevices(c echo.Context) error {
	slog.Info("Render devices section for the home page")

	// Get devices...
	devices, err := h.registry.Devices.List()
	if err != nil {
		return errors.Wrap(err, "failed to list devices")
	}

	rDevices, err := services.ResolveDevices(h.registry, devices...)
	if err != nil {
		return errors.Wrap(err, "failed to resolve devices")
	}

	return components.SectionDevices(false, rDevices).Render(c.Request().Context(), c.Response())
}

func (h *Handler) GetGroups(c echo.Context) error {
	slog.Info("Render groups section for the home page")

	// Get groups...
	groups, err := h.registry.Groups.List()
	if err != nil {
		return errors.Wrap(err, "failed to list groups")
	}

	// ...resolve them
	resolvedGroups, err := services.ResolveGroups(h.registry, groups...)
	if err != nil {
		return errors.Wrap(err, "failed to resolve groups")
	}

	return components.SectionGroups(false, resolvedGroups).Render(c.Request().Context(), c.Response())
}

func (h *Handler) DeleteDevice(c echo.Context) error {
	deviceID, err := utils.QueryParamDeviceID(c, "id", false)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "failed to get device ID from query parameter"))
	}

	slog.Info("Delete a device", "id", deviceID)

	if err = h.registry.Devices.Delete(deviceID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "failed to delete device"))
	}

	c.Response().Header().Set("HX-Trigger", "reloadDevices")
	return nil
}

func (h *Handler) DeleteGroup(c echo.Context) error {
	groupID, err := utils.QueryParamGroupID(c, "id", false)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "failed to get group ID from query parameter"))
	}

	slog.Info("Delete a group", "id", groupID)

	if err = h.registry.Groups.Delete(groupID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "failed to delete group"))
	}

	c.Response().Header().Set("HX-Trigger", "reloadGroups")
	return nil
}

func (h *Handler) PostTogglePowerDevice(c echo.Context) error {
	deviceID, err := utils.QueryParamDeviceID(c, "id", false)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "failed to get device ID from query parameter"))
	}

	slog.Info("Toggle power for a device", "id", deviceID)

	color, err := h.registry.DeviceControls.TogglePower(deviceID)
	if err != nil {
		err = errors.Wrap(err, "failed to toggle power for device %d", deviceID)
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
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "failed to get group ID from query parameter"))
	}

	group, err := h.registry.Groups.Get(groupID)
	if err != nil {
		oob.OOBRenderPageHomeGroupError(c, groupID, []error{err})

		if services.IsNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, errors.Wrap(err, "group not found"))
		}

		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "failed to get group"))
	}

	slog.Info("Turn on a group", "id", groupID, "devices", group.Devices)

	devices, err := h.registry.Devices.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "failed to list devices"))
	}

	wg := &sync.WaitGroup{}
	errs := make([]error, 0)
	for _, d := range devices {
		wg.Go(func() {
			if !slices.Contains(group.Devices, d.ID) {
				if err := h.registry.DeviceControls.TurnOff(d.ID); err != nil {
					if device, err2 := h.registry.Devices.Get(d.ID); err2 != nil {
						errs = append(errs, errors.Wrap(err2, "failed to get device %d from the database", d.ID))
					} else {
						errs = append(errs, errors.Wrap(err, "failed to turn off device \"%s\", which is not in this group", device.Name))
					}
				}

				return
			}

			if err := h.registry.DeviceControls.TurnOn(d.ID); err != nil {
				if device, err2 := h.registry.Devices.Get(d.ID); err2 != nil {
					errs = append(errs, errors.Wrap(err2, "failed to get device %d from the database", d.ID))
				} else {
					errs = append(errs, errors.Wrap(err, "failed to turn on device \"%s\"", device.Name))
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
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "failed to get group ID from query parameter"))
	}

	group, err := h.registry.Groups.Get(groupID)
	if err != nil {
		oob.OOBRenderPageHomeGroupError(c, groupID, []error{err})

		if services.IsNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, errors.Wrap(err, "group not found"))
		}

		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "failed to get group"))
	}

	slog.Info("Turn off a group", "id", groupID, "devices", group.Devices)

	devices, err := h.registry.Devices.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "failed to list devices"))
	}

	wg := &sync.WaitGroup{}
	errs := make([]error, 0)
	for _, d := range devices {
		wg.Go(func() {
			if err := h.registry.DeviceControls.TurnOff(d.ID); err != nil {
				if device, err2 := h.registry.Devices.Get(d.ID); err2 != nil {
					errs = append(errs, errors.Wrap(err2, "failed to get device %d from the database", d.ID))
				} else {
					errs = append(errs, errors.Wrap(err, "failed to turn off device \"%s\"", device.Name))
				}
			}
		})
	}
	wg.Wait()

	oob.OOBRenderPageHomeGroupError(c, groupID, errs)

	return nil
}
