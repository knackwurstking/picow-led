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

	apiDevices(e, o.Api)
	frontend(e, o.Frontend)
	pwa(e, o.PWA)
}
