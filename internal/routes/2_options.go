package routes

import (
	"io/fs"
	"picow-led/internal/api"
)

type Options struct {
	Global
	Templates fs.FS
	Config    *api.Config
}

func (o *Options) Devices() Devices {
	g := o.Global
	g.SubTitle = "Devices"

	return Devices{
		Global: g,
	}
}
