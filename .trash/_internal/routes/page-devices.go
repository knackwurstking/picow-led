package routes

import "picow-led/internal/types"

type Devices struct {
	Global
}

func (d Devices) List() []*types.Device {
	return cache.Devices()
}
