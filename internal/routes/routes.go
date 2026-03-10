package routes

import (
	"github.com/knackwurstking/picow-led/internal/assets"
	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/internal/handlers"
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo) {
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
}
