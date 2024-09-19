package main

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"

	_api "github.com/knackwurstking/picow-led-server/pkg/api"
	_clients "github.com/knackwurstking/picow-led-server/pkg/clients"
)

var (
	api                = _api.NewAPI()
	e       *echo.Echo = echo.New()
	clients            = _clients.NewClients()
	config             = NewConfig()
	flags              = NewFlags(Port)
	wg                 = &sync.WaitGroup{}
)

func main() {
	flags.Parse()

	// Parse version and exit

	if flags.Version {
		fmt.Printf("picow-led-server %s\n", Version)
		os.Exit(ErrorCodeOK)
	}

	// Handle config

	slog.Debug("Try to load API configuration", "config.Path", config.Path)
	if err := config.load(); err != nil {
		slog.Error("Loading API configuration", "config.Path", config.Path, "err", err)
		os.Exit(ErrorCodeConfiguration)
	}

	// Set middleware

	e.Use(
		slogecho.New(slog.Default()),
		middleware.Recover(),
		middleware.CORS(),
	)

	// Create endpoints

	endpointsStatic()
	endpointsEvents()
	endpointsApi()

	// Echo server configuration

	e.HideBanner = true
	e.HidePort = true

	// Start server

	slog.Info("HTTP server started", "Host", Host, "Port", flags.Port)
	addr := fmt.Sprintf("%s:%d", Host, flags.Port)
	if err := e.Start(addr); err != nil {
		slog.Error("Server start", "err", err)
		os.Exit(ErrorCodeServerStart)
	}
}
