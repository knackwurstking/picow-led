package routes

import (
	"picow-led/internal/api"

	"github.com/labstack/echo/v4"
)

type Options struct {
	ServerPathPrefix string      `json:"server-path-prefix"`
	Version          string      `json:"version"`
	Api              *api.Config `json:"api"`
}

func Create(e *echo.Echo, data Options) {
	FrontendCache = api.GetDevices(data.Api)

	apiDevices(e, data)
	frontend(e, data)
	pwa(e, data)
}
