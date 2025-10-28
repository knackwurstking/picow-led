package main

import (
	picowled "github.com/knackwurstking/picow-led"
	"github.com/knackwurstking/picow-led/env"
	"github.com/knackwurstking/picow-led/handlers"
	"github.com/knackwurstking/picow-led/services"
	"github.com/labstack/echo/v4"
)

func router(e *echo.Echo, r *services.Registry) {
	e.StaticFS(env.Args.ServerPathPrefix+"/", picowled.GetAssets())

	handlers.NewPages(r).Register(e)
	handlers.NewHXHome(r).Register(e)
}
