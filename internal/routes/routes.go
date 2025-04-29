package routes

import (
	"picow-led/internal/api"

	"github.com/labstack/echo/v4"
)

type Options struct {
	Api      Api
	Frontend Frontend
	PWA      PWA
}

func Create(e *echo.Echo, o Options) {
	FrontendCache = api.GetDevices(o.Api.Config)

	apiDevices(e, Api{
		ServerPathPrefix: o.Api.ServerPathPrefix,
		Config:           o.Api.Config,
	})

	frontend(e, Frontend{
		ServerPathPrefix: o.Frontend.ServerPathPrefix,
	})

	pwa(e, PWA{
		ServerPathPrefix: o.PWA.ServerPathPrefix,
		Version:          o.PWA.Version,
	})
}
