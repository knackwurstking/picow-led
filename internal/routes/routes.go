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

	// API Endpoints
	group := e.Group(env.Route("/api"))
	{
		// Register API endpoints here
	}

	// UI Endpoints
	group = e.Group(env.Route(""))
	{
		// Register UI endpoints here
		group.GET("/", handlers.Home)
	}

	// HTMX Endpoints
	group = e.Group(env.Route("/htmx"))
	{
		// Register HTMX endpoints here
		group.GET("/devices", handlers.HTMXDevices(r))
		subGroup := group.Group("/dialogs")
		{
			subGroup.GET("/add-device", handlers.HTMXAddDeviceDialog(r, http.MethodGet))
			subGroup.POST("/add-device", handlers.HTMXAddDeviceDialog(r, http.MethodPost))
			subGroup.GET("/edit-device", handlers.HTMXAddDeviceDialog(r, http.MethodGet))
			subGroup.POST("/edit-device", handlers.HTMXAddDeviceDialog(r, http.MethodPost))
		}

	}
}
