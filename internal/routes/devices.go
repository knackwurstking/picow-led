package routes

import "picow-led/internal/api"

type Devices struct {
	Global
}

func (d Devices) List() []*api.Device {
	return cache.Devices()
}
