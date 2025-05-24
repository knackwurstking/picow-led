package routes

import (
	"io/fs"
	"picow-led/internal/api"
	"picow-led/internal/types"
)

type Options struct {
	Global
	Templates fs.FS
	Config    *types.APIConfig
	WS        *api.WS
}

func (o *Options) Devices() Devices {
	g := o.Global
	g.SubTitle = "Devices"

	return Devices{
		Global: g,
	}
}

func (o *Options) PageControl(addr string) PageControl {
	// TODO: Get device for this address
	device :=

	g := o.Global
	g.SubTitle = addr // TODO: Set name if available

	return PageControl{
		Global: g,
		Device: device,
	}
}
