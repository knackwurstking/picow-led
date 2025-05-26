package routes

import (
	"picow-led/internal/database"

	"github.com/labstack/echo/v4"
)

type Options struct {
	ServerPathPrefix string
	DB               *database.DB
}

func Register(e *echo.Echo, o *Options) {
	apiHandler := NewAPIHandler(o.DB)

	apiGroup := e.Group(o.ServerPathPrefix + "/api")

	apiGroup.GET("/devices", apiHandler.GetDevices)

	apiGroup.GET("/devices/:addr", apiHandler.GetDevicesAddr)
	apiGroup.GET("/devices/:addr/name", apiHandler.GetDevicesAddrName)
	apiGroup.GET("/devices/:addr/active_color", apiHandler.GetDevicesAddrColor)
	apiGroup.GET("/devices/:addr/color", apiHandler.GetDevicesAddrColor)
	apiGroup.GET("/devices/:addr/pins", apiHandler.GetDevicesAddrPins)

	apiGroup.GET("/devices/:addr/power", apiHandler.GetDevicesAddrPower)
	apiGroup.POST("/devices/:addr/power", apiHandler.PostDevicesAddrPower)

	apiGroup.GET("/colors", apiHandler.GetColors)
	apiGroup.POST("/colors", apiHandler.PostColors)
	apiGroup.PUT("/colors", apiHandler.PutColors)

	apiGroup.GET("/colors/:id", apiHandler.GetColorsID)
	apiGroup.POST("/colors/:id", apiHandler.PostColorsID)
	apiGroup.DELETE("/colors/:id", apiHandler.DeleteColorsID)
}
