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
	if cache != nil {
		cache.Close()
	}

	ws := api.NewWS(e.Logger)

	cache = NewCache(ws)
	cache.SetDevices(api.GetDevices(o.Config)...)

	go ws.Start()

	apiRoutes(e, apiOptions{
		ServerPathPrefix: o.ServerPathPrefix,
		Config:           o.Config,
	})

	wsRoutes(e, wsOptions{
		ServerPathPrefix: o.ServerPathPrefix,
		WS:               ws,
	})

	frontendRoutes(e, frontendOptions{
		ServerPathPrefix: o.ServerPathPrefix,
		Version:          o.Version,
		Templates:        o.Templates,
	})
}
