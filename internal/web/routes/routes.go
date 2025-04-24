package routes

import (
	"picow-led/internal/web/api"

	"github.com/labstack/echo/v4"
)

type Data struct {
	ServerPathPrefix string
	Api              api.Data
}

func Create(e *echo.Echo, data Data) {
	apiDevices(e, data)
}
