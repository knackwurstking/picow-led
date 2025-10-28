package handlers

import (
	"net/http"

	"github.com/knackwurstking/picow-led/services"
	"github.com/labstack/echo/v4"
)

type HXDialogs struct {
	registry *services.Registry
}

func NewHXDialogs(r *services.Registry) *HXDialogs {
	return &HXDialogs{
		registry: r,
	}
}

func (h HXDialogs) Register(e *echo.Echo) {
	Register(e, http.MethodGet, "/htmx/dialog/edit-device", h.EditDevice)
	Register(e, http.MethodGet, "/htmx/dialog/edit-group", h.EditGroup)
}

func (h HXDialogs) EditDevice(c echo.Context) error {
	// TODO: ...

	return nil
}

func (h HXDialogs) EditGroup(c echo.Context) error {
	// TODO: ...

	return nil
}
