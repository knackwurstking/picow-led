package routes

import (
	"io/fs"
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
	ServerPathPrefix string
	Version          string
	Templates        fs.FS
	Config           *api.Config
}

type Cache struct {
	// TODO: Handle the mutex lock internal, add methods for devices and
	// 		 color, make fields private
	Devices []*api.Device
	Color   []api.MicroColor
	Mutex   *sync.Mutex
}

func Create(e *echo.Echo, o Options) {
	cache.Devices = api.GetDevices(o.Config)

	apiRoutes(e, Api{
		ServerPathPrefix: o.ServerPathPrefix,
		Config:           o.Config,
	})
	wsRoutes(e, WS{
		ServerPathPrefix: o.ServerPathPrefix,
	})
	frontendRoutes(e, Frontend{
		ServerPathPrefix: o.ServerPathPrefix,
		Version:          o.Version,
		Templates:        o.Templates,
	})
}
