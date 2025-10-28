package main

import (
	"github.com/knackwurstking/picow-led/handlers"
	"github.com/knackwurstking/picow-led/services"
	"github.com/labstack/echo/v4"
)

func router(e *echo.Echo, r *services.Registry) {
	// TODO: Server static files first

	handlers.NewPages(r).Register(e)
	handlers.NewHXHome(r).Register(e)
}
