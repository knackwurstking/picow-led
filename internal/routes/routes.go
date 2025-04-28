package routes

import (
	"picow-led/internal/api"

	"github.com/labstack/echo/v4"
)

type Options struct {
	ServerPathPrefix string      `json:"server-path-prefix" yaml:"server-path-prefix"`
	Api              *api.Config `json:"api" yaml:"api"`
}

func Create(e *echo.Echo, data Options) {
	FrontendCache = api.GetDevices(data.Api)

	apiDevices(e, data)
	frontend(e, data)
}
