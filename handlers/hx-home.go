package handlers

import (
	"net/http"

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
	// TODO: Continue here

	return nil
}

func (h HXHome) SectionGroups(c echo.Context) error {
	// TODO: Continue here

	return nil
}
