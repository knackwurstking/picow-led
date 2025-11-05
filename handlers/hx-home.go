package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/knackwurstking/picow-led/components"
	"github.com/knackwurstking/picow-led/services"
	"github.com/labstack/echo/v4"
)

type HXHome struct {
	registry *services.Registry
}

func NewHxHome(r *services.Registry) *HXHome {
	return &HXHome{
		registry: r,
	}
}

func (h HXHome) Register(e *echo.Echo) {
	Register(e, http.MethodGet, "/htmx/home/section/devices", h.GetSectionDevices)
	Register(e, http.MethodGet, "/htmx/home/section/groups", h.GetSectionGroups)
}

func (h HXHome) GetSectionDevices(c echo.Context) error {
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

	return components.PageHome_SectionDevices(
		false, rDevices...,
	).Render(c.Request().Context(), c.Response())
}

func (h HXHome) GetSectionGroups(c echo.Context) error {
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

	return components.PageHome_SectionGroups(
		false,
		resolvedGroups...,
	).Render(c.Request().Context(), c.Response())
}
