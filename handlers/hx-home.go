package handlers

import (
	"github.com/knackwurstking/picow-led/env"
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
	e.GET(env.Args.ServerPathPrefix+"/htmx/home/section/devices", h.SectionDevices)
	e.GET(env.Args.ServerPathPrefix+"/htmx/home/section/groups", h.SectionGroups)
}

func (h HXHome) SectionDevices(c echo.Context) error {
	// TODO: ...

	return nil
}

func (h HXHome) SectionGroups(c echo.Context) error {
	// TODO: ...

	return nil
}
