package routes

import (
	"fmt"
	"picow-led/internal/api"
	"sync"
)

var cache = NewCache()

type Cache struct {
	devices []*api.Device
	color   []api.MicroColor

	mutex *sync.Mutex
}

func NewCache() *Cache {
	return &Cache{
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

func (c *Cache) Devices() []*api.Device {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.devices
}

func (c *Cache) SetDevices(devices ...*api.Device) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.devices = devices
}

func (c *Cache) UpdateDevice(addr string, device *api.Device) (*api.Device, error) {
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

func (c *Cache) Color() []api.MicroColor {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.color
}

func (c *Cache) UpdateColor(index int, color api.MicroColor) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.color)-1 < index {
		return fmt.Errorf("no such index: %d", index)
	}

	c.color[index] = color

	return nil
}
