package routes

import (
	"picow-led/internal/api"

	"github.com/labstack/echo/v4"
)

type Options struct {
	ServerPathPrefix string      `json,yaml:"server-path-prefix"`
	Version          string      `json,yaml:"version"`
	Api              *api.Config `json,yaml:"api"`
}

func Create(e *echo.Echo, data Options) {
	apiDevices(e, data)
	frontend(e, data)
}
