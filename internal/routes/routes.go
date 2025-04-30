package routes

import (
	"picow-led/internal/api"

	"github.com/labstack/echo/v4"
)

var cache = &Cache{
	Devices: make([]*api.Device, 0),
	Color:   make(map[string]api.MicroColor),
}

type Options struct {
	Api      Api
	Frontend Frontend
}

type Cache struct {
	Devices []*api.Device
	Color   map[string]api.MicroColor
}

func Create(e *echo.Echo, o Options) {
	cache.Devices = api.GetDevices(o.Api.Config)

	apiDevices(e, o.Api)
	frontend(e, o.Frontend)
}
