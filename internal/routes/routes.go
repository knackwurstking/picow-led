package routes

import (
	"io/fs"
	"picow-led/internal/api"

	"github.com/labstack/echo/v4"
)

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
