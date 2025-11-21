package groups

import (
	"log/slog"
	"net/http"
	"slices"
	"sync"

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
		"/htmx/home/groups", h.GetGroups)
	utils.Register(e, http.MethodDelete,
		"/htmx/home/groups/delete", h.DeleteGroup)
	utils.Register(e, http.MethodPost,
		"/htmx/home/groups/turn-on", h.PostTurnOnGroup)
	utils.Register(e, http.MethodPost,
		"/htmx/home/groups/turn-off", h.PostTurnOffGroup)
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

	return components.SectionGroups(false, resolvedGroups).Render(c.Request().Context(), c.Response())
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
