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

	apiGroup.GET(o.ServerPathPrefix+"/api/devices", apiHandler.GetDevices)

	apiGroup.GET(o.ServerPathPrefix+"/api/devices/:addr",
		apiHandler.GetDevicesAddr)

	apiGroup.GET(o.ServerPathPrefix+"/api/devices/:addr/name",
		apiHandler.GetDevicesAddrName)

	apiGroup.GET(o.ServerPathPrefix+"/api/devices/:addr/active_color",
		apiHandler.GetDevicesAddrColor)

	apiGroup.GET(o.ServerPathPrefix+"/api/devices/:addr/color",
		apiHandler.GetDevicesAddrColor)

	apiGroup.GET(o.ServerPathPrefix+"/api/devices/:addr/pins",
		apiHandler.GetDevicesAddrPins)

	apiGroup.GET(o.ServerPathPrefix+"/api/devices/:addr/power",
		apiHandler.GetDevicesAddrPower)
	apiGroup.POST(o.ServerPathPrefix+"/api/devices/:addr/power",
		apiHandler.PostDevicesAddrPower)

	apiGroup.GET(o.ServerPathPrefix+"/api/colors", apiHandler.GetColors)
	apiGroup.POST(o.ServerPathPrefix+"/api/colors", apiHandler.PostColors)
	apiGroup.PUT(o.ServerPathPrefix+"/api/colors", apiHandler.PutColors)

	apiGroup.GET(o.ServerPathPrefix+"/api/colors/:index",
		apiHandler.GetColorsIndex)
	apiGroup.PUT(o.ServerPathPrefix+"/api/colors/:index",
		apiHandler.PutColorsIndex)
	apiGroup.DELETE(o.ServerPathPrefix+"/api/colors/:index",
		apiHandler.DeleteColorsIndex)
}
