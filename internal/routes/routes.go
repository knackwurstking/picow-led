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

	apiGroup.GET("/devices/:addr", apiHandler.GetDevice)

	apiGroup.GET("/devices/:addr/name", apiHandler.GetDeviceName)

	apiGroup.GET("/devices/:addr/active_color", apiHandler.GetDeviceColor)

	apiGroup.GET("/devices/:addr/color", apiHandler.GetDeviceColor)
	apiGroup.POST("/devices/:addr/color", apiHandler.PostDeviceColor)

	apiGroup.GET("/devices/:addr/pins", apiHandler.GetDevicePins)

	apiGroup.GET("/devices/:addr/power", apiHandler.GetDevicePower)
	apiGroup.POST("/devices/:addr/power", apiHandler.PostDevicePower)

	apiGroup.GET("/colors", apiHandler.GetColors)
	apiGroup.POST("/colors", apiHandler.PostColors)
	apiGroup.PUT("/colors", apiHandler.PutColors)

	apiGroup.GET("/colors/:id", apiHandler.GetColorsID)
	apiGroup.POST("/colors/:id", apiHandler.PostColorsID)
	apiGroup.DELETE("/colors/:id", apiHandler.DeleteColorsID)
}
