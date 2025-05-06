package routes

import (
	"fmt"
	"io/fs"
	"picow-led/internal/api"
	"sync"

	"github.com/labstack/echo/v4"
)

var cache = newCache()

type _cache struct {
	devices []*api.Device
	color   []api.MicroColor

	mutex *sync.Mutex
}

func newCache() *_cache {
	return &_cache{
		devices: make([]*api.Device, 0),
		color: []api.MicroColor{
			{255, 255, 255, 255},
			{255, 0, 0, 0},
			{0, 255, 0, 0},
			{0, 0, 255, 0},
		},

		mutex: &sync.Mutex{},
	}
}

func (c *_cache) Devices() []*api.Device {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.devices
}

func (c *_cache) SetDevices(devices ...*api.Device) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.devices = devices
}

func (c *_cache) UpdateDevice(addr string, device *api.Device) (*api.Device, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, d := range c.devices {
		if d.Server.Addr != addr {
			continue
		}

		// Only merge things changed after PostDevicesColor call
		d.Color = device.Color
		d.Error = device.Error
		d.Online = device.Online

		// Data to return
		return d, nil
	}

	return nil, fmt.Errorf("device \"%s\" not found", addr)
}

func (c *_cache) Color() []api.MicroColor {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.color
}

func (c *_cache) UpdateColor(index int, color api.MicroColor) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.color)-1 < index {
		return fmt.Errorf("no such index: %d", index)
	}

	c.color[index] = color

	return nil
}

type Options struct {
	ServerPathPrefix string
	Version          string
	Templates        fs.FS
	Config           *api.Config
}

func Create(e *echo.Echo, o Options) {
	cache.SetDevices(api.GetDevices(o.Config)...)

	apiRoutes(e, apiOptions{
		ServerPathPrefix: o.ServerPathPrefix,
		Config:           o.Config,
	})

	o.Config.WS.Start()
	wsRoutes(e, wsOptions{
		ServerPathPrefix: o.ServerPathPrefix,
		WS:               o.Config.WS,
	})

	frontendRoutes(e, frontendOptions{
		ServerPathPrefix: o.ServerPathPrefix,
		Version:          o.Version,
		Templates:        o.Templates,
	})
}
