package handlers

import (
	"net/http"

	"github.com/knackwurstking/picow-led/components"
	"github.com/knackwurstking/picow-led/services"
	"github.com/labstack/echo/v4"
)

type HXHome struct {
	registry *services.Registry
}

func NewHXHome(r *services.Registry) *HXHome {
	return &HXHome{
		registry: r,
	}
}

func (h HXHome) Register(e *echo.Echo) {
	Register(e, http.MethodGet, "/htmx/home/section/devices", h.SectionDevices)
	Register(e, http.MethodGet, "/htmx/home/section/groups", h.SectionGroups)
}

func (h HXHome) SectionDevices(c echo.Context) error {
	// Get devices...
	devices, err := h.registry.Devices.List()
	if err != nil {
		return err
	}

	// ...resolve them
	resolvedDevices, err := services.ResolveDevices(h.registry, devices...)
	if err != nil {
		return err
	}

	return components.PageHome_SectionDevices(
		false,
		resolvedDevices...,
	).Render(c.Request().Context(), c.Response())
}

func (h HXHome) SectionGroups(c echo.Context) error {
	// Get groups...
	groups, err := h.registry.Groups.List()
	if err != nil {
		return err
	}

	// ...resolve them
	resolvedGroups, err := services.ResolveGroups(h.registry, groups...)
	if err != nil {
		return err
	}

	return components.PageHome_SectionGroups(
		false,
		resolvedGroups...,
	).Render(c.Request().Context(), c.Response())
}
