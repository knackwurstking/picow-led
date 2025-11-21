package handlers

import (
	"github.com/knackwurstking/picow-led/handlers/devices"
	"github.com/knackwurstking/picow-led/handlers/dialogs"
	"github.com/knackwurstking/picow-led/handlers/groups"
	"github.com/knackwurstking/picow-led/services"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Register(e *echo.Echo)
}

func GetAll(r *services.Registry) []Handler {
	return []Handler{
		devices.NewHandler(r),
		groups.NewHandler(r),
		dialogs.NewHandler(r),
	}
}
