package routes

import (
	"io/fs"
	"picow-led/internal/types"
)

type Options struct {
	Global
	Templates fs.FS
	Config    *types.APIConfig
}

func (o *Options) Devices() Devices {
	g := o.Global
	g.SubTitle = "Devices"

	return Devices{
		Global: g,
	}
}
