package routes

import (
	"net/http"

	"github.com/knackwurstking/picow-led/internal/assets"
	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/internal/handlers"
	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, r *services.Registry) {
	// Static Assets
	assets.ServePublicFS(e)

	{ // Register API endpoints here
		group := e.Group(env.Route("/api"))
		{
			devicesGroup := group.Group("/devices/:id")

			devicesGroup.POST("/color", handlers.APISetDeviceColor(r, http.MethodPost))
			devicesGroup.POST("/white", handlers.APISetDeviceWhite(r, http.MethodPost))
			devicesGroup.POST("/white2", handlers.APISetDeviceWhite2(r, http.MethodPost))
			devicesGroup.POST("/brightness", handlers.APISetDeviceBrightness(r, http.MethodPost))

			devicesGroup.POST("/rgbw", handlers.APISetDeviceRGBW(r, http.MethodPost))
		}
	}

	{ // Register UI endpoints here
		group := e.Group(env.Route(""))

		group.GET("/", handlers.Home)
		group.GET("/device", handlers.Device(r, http.MethodGet))
	}

	{ // Register HTMX endpoints here
		group := e.Group(env.Route("/htmx"))

		group.GET("/devices", handlers.HTMXDevices(r))
		group.POST("/devices/toggle-power", handlers.HTMXToggleDevicePower(r))

		{
			dialogsGroup := group.Group("/dialogs")

			dialogsGroup.GET("/add-device", handlers.HTMXAddDeviceDialog(r, http.MethodGet))
			dialogsGroup.POST("/add-device", handlers.HTMXAddDeviceDialog(r, http.MethodPost))

			dialogsGroup.GET("/edit-device", handlers.HTMXEditDeviceDialog(r, http.MethodGet))
			dialogsGroup.POST("/edit-device", handlers.HTMXEditDeviceDialog(r, http.MethodPost))
			dialogsGroup.DELETE("/edit-device", handlers.HTMXEditDeviceDialog(r, http.MethodDelete))
		}
	}
}
