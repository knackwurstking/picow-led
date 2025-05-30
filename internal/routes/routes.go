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
	registerGroup(e.Group(o.ServerPathPrefix+"/api"), apiHandler)

	// TODO: Register /ws routes
}

func registerGroup(g *echo.Group, apiHandler *APIHandler) {
	g.GET("/devices", apiHandler.GetDevices)

	g.GET("/devices/:addr", apiHandler.GetDevice)

	g.GET("/devices/:addr/name", apiHandler.GetDeviceName)

	g.GET("/devices/:addr/active_color", apiHandler.GetDeviceActiveColor)

	g.GET("/devices/:addr/color", apiHandler.GetDeviceColor)
	g.POST("/devices/:addr/color", apiHandler.PostDeviceColor)

	g.GET("/devices/:addr/pins", apiHandler.GetDevicePins)

	g.GET("/devices/:addr/power", apiHandler.GetDevicePower)
	g.POST("/devices/:addr/power", apiHandler.PostDevicePower)

	g.GET("/colors", apiHandler.GetColors)
	g.POST("/colors", apiHandler.PostColors)
	g.PUT("/colors", apiHandler.PutColors)

	g.GET("/colors/:id", apiHandler.GetColorsID)
	g.POST("/colors/:id", apiHandler.PostColorsID)
	g.DELETE("/colors/:id", apiHandler.DeleteColorsID)
}
