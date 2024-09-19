package main

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"

	"github.com/knackwurstking/picow-led-server/cmd/picow-led-server/endpoints"
)

var (
	e     *echo.Echo = echo.New()
	flags            = NewFlags(Port)
	wg               = &sync.WaitGroup{}
)

func main() {
	flags.Parse()

	// Parse version and exit

	if flags.Version {
		fmt.Printf("picow-led-server %s\n", Version)
		os.Exit(ErrorCodeOK)
	}

	// Set middleware

	e.Use(
		slogecho.New(slog.Default()),
		middleware.Recover(),
		middleware.CORS(),
	)

	// Echo server configuration

	e.HideBanner = true
	e.HidePort = true

	// Create endpoints

	endpoints.Create(e)

	// Start server

	slog.Info("HTTP server started", "Host", Host, "Port", flags.Port)
	addr := fmt.Sprintf("%s:%d", Host, flags.Port)
	if err := e.Start(addr); err != nil {
		slog.Error("Server start", "err", err)
		os.Exit(ErrorCodeServerStart)
	}
}
