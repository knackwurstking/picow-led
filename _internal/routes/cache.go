package routes

import (
	"fmt"
	"picow-led/internal/api"
	"picow-led/internal/types"
	"slices"
	"sync"
)

var cache *Cache

type Cache struct {
	devices []*types.Device
	colors  []types.MicroColor

	mutex *sync.Mutex
	ws    *api.WS
}

func NewCache(ws *api.WS) *Cache {
	return &Cache{
		devices: make([]*types.Device, 0),
		colors: []types.MicroColor{
			{255, 255, 255, 255},
			{255, 0, 0, 0},
			{0, 255, 0, 0},
			{0, 0, 255, 0},
		},

		mutex: &sync.Mutex{},
		ws:    ws,
	}
}

func (c *Cache) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.ws != nil {
		c.ws.Stop()
	}
}

func (c *Cache) Devices() []*types.Device {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.devices
}

func (c *Cache) Device(addr string) *types.Device {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, d := range c.devices {
		if d.Server.Addr == addr {
			return d
		}
	}

	return nil
}

func (c *Cache) SetDevices(devices ...*types.Device) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.devices = devices
}

func (c *Cache) UpdateDevice(addr string, device *types.Device) (*types.Device, error) {
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

		if c.ws != nil {
			c.ws.BroadcastDevice(d)
		}

		return d, nil
	}

	return nil, fmt.Errorf("device \"%s\" not found", addr)
}

func (c *Cache) Colors() []types.MicroColor {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.colors
}

func (c *Cache) UpdateColor(index int, color types.MicroColor) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.colors)-1 < index {
		return fmt.Errorf("no such index: %d", index)
	}

	if index < 0 {
		c.colors = append(c.colors, color)
	} else {
		c.colors[index] = color
	}

	if c.ws != nil {
		c.ws.BroadcastColors(c.colors)
	}

	return nil
}

func (c *Cache) DeleteColor(index int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.colors = slices.Delete(c.colors, index, index+1)

	if c.ws != nil {
		go c.ws.BroadcastColors(c.colors)
	}

	return nil
}
