package main

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/knackwurstking/ui"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/internal/routes"
	"github.com/knackwurstking/picow-led/internal/services"

	_ "github.com/mattn/go-sqlite3"
)

const (
	ExitCodeSuccess = iota
	ExitCodeGeneric
	ExitCodeServerStart
	ExitCodeDatabase
	ExitCodeRegistry
)

var (
	log *ui.Logger = env.NewLogger("main")
)

func main() {
	// Open SQL database and pass it to the registry
	db, err := sql.Open("sqlite3", filepath.Join(env.DBPath, "picow-led.sqlite"))
	if err != nil {
		log.Error("Failed to open database: %v", err)
		os.Exit(ExitCodeDatabase)
	}

	db.SetMaxOpenConns(25)   // Allow up to 25 open connections
	db.SetMaxIdleConns(25)   // Allow up to 25 idle connections
	db.SetConnMaxLifetime(0) // No maximum lifetime

	registry, err := services.NewRegistry(db)
	if err != nil {
		log.Error("Failed to initialize registry: %v", err)
		os.Exit(ExitCodeRegistry)
	}

	// Start the server
	e := echo.New()

	// Middleware
	requestLogger := env.NewLogger("request")
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:  true,
		LogStatus:  true,
		LogURI:     true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			// Example:
			//value, _ := c.Get("context-value").(int)
			if v.Error != nil {
				requestLogger.Error("%v", v.Error)
			}

			requestLogger.Info("%s %d %s %s", v.Method, v.Status, v.URI, v.Latency)

			return nil
		},
	}))
	e.Use(middleware.AddTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(ui.EchoMiddlewareCache())

	// Register handlers
	routes.Register(e, registry)

	if err := e.Start(env.ServerAddress); err != nil {
		log.Error("Failed to start server: %v", err)
		os.Exit(ExitCodeServerStart)
	}
}
