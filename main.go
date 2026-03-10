package main

import (
	"os"

	"github.com/knackwurstking/ui"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/internal/routes"
)

const (
	ExitCodeServerStart = 2
)

var (
	log *ui.Logger = env.NewLogger("main")
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.AddTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(ui.EchoMiddlewareCache())

	routes.Register(e)

	if err := e.Start(env.ServerAddress); err != nil {
		log.Error("Failed to start server: %v", err)
		os.Exit(ExitCodeServerStart)
	}
}
