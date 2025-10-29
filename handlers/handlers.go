package handlers

import (
	"github.com/knackwurstking/picow-led/services"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Register(e *echo.Echo)
}

func GetAll(r *services.Registry) []Handler {
	return []Handler{
		NewPages(r),
		NewHxHome(r),
		NewHxDialogs(r),
		NewHxDevices(r),
	}
}
