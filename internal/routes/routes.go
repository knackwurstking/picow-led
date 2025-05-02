package routes

import (
	"picow-led/internal/api"

	"github.com/labstack/echo/v4"
)

var cache = &Cache{
	Devices: make([]*api.Device, 0),
	Color: []api.MicroColor{
		{255, 255, 255, 255},
		{255, 0, 0, 0},
		{0, 255, 0, 0},
		{0, 0, 255, 0},
	},
}

type Options struct {
	Api      Api
	Frontend Frontend
}

type Cache struct {
	Devices []*api.Device
	Color   []api.MicroColor
}

func Create(e *echo.Echo, o Options) {
	cache.Devices = api.GetDevices(o.Api.Config)

	apiRoutes(e, o.Api)
	frontendRoutes(e, o.Frontend)
}
