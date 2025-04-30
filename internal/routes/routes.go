package routes

import (
	"picow-led/internal/api"

	"github.com/labstack/echo/v4"
)

type Options struct {
	Api      Api
	Frontend Frontend
}

func Create(e *echo.Echo, o Options) {
	cache.Devices = api.GetDevices(o.Api.Config)

	apiDevices(e, o.Api)
	frontend(e, o.Frontend)
}
