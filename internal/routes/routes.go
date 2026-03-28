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

		group.GET("/groups", handlers.HTMXGroups(r))
		group.POST("/groups/power", handlers.HTMXPowerGroup(r))

		group.GET("/device/pins", handlers.HTMXDevicePins(r, http.MethodGet))
		group.POST("/device/pins", handlers.HTMXDevicePins(r, http.MethodPost))

		{
			dialogsGroup := group.Group("/dialogs")

			// Add Device
			dialogsGroup.GET("/add-device", handlers.HTMXAddDeviceDialog(r, http.MethodGet))
			dialogsGroup.POST("/add-device", handlers.HTMXAddDeviceDialog(r, http.MethodPost))

			// Edit Device
			dialogsGroup.GET("/edit-device", handlers.HTMXEditDeviceDialog(r, http.MethodGet))
			dialogsGroup.POST("/edit-device", handlers.HTMXEditDeviceDialog(r, http.MethodPost))
			dialogsGroup.DELETE("/edit-device", handlers.HTMXEditDeviceDialog(r, http.MethodDelete))

			// Add Group
			dialogsGroup.GET("/add-group", handlers.HTMXAddGroupDialog(r, http.MethodGet))
			dialogsGroup.POST("/add-group", handlers.HTMXAddGroupDialog(r, http.MethodPost))

			// Edit Group
			dialogsGroup.GET("/edit-group", handlers.HTMXEditGroupDialog(r, http.MethodGet))
			dialogsGroup.POST("/edit-group", handlers.HTMXEditGroupDialog(r, http.MethodPost))
			dialogsGroup.DELETE("/edit-group", handlers.HTMXEditGroupDialog(r, http.MethodDelete))
		}
	}
}
