package routes

import (
	"picow-led/internal/api"
	"sync"

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
	Mutex: &sync.Mutex{},
}

type Options struct {
	Api      Api
	Frontend Frontend
}

type Cache struct {
	// TODO: Handle the mutex lock internal, add methods for devices and
	// 		 color, make fields private
	Devices []*api.Device
	Color   []api.MicroColor
	Mutex   *sync.Mutex
}

func Create(e *echo.Echo, o Options) {
	cache.Devices = api.GetDevices(o.Api.Config)

	apiRoutes(e, o.Api)
	frontendRoutes(e, o.Frontend)
}
