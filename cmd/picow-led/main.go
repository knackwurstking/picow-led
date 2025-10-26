package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/knackwurstking/picow-led/env"
	"github.com/labstack/echo/v4"
)

func main() {
	parseFlags()
	initializeLogging()
	startServer()
}

func parseFlags() {
	var logFormat string = string(env.Args.LogFormat)

	// Server Address
	flag.StringVar(&env.Args.Addr, "addr", env.Args.Addr, "Server address")

	// Server Path Prefix
	flag.StringVar(&env.Args.ServerPathPrefix, "path-prefix", env.Args.ServerPathPrefix, "Server path prefix")

	// Debug Flag
	flag.BoolVar(&env.Args.Debug, "debug", env.Args.Debug, "Enable debug mode")

	// Log Format: "text", "json"
	flag.StringVar(&logFormat, "log-format", logFormat, "Log format: text, json")

	// Custom Usage Message
	flag.Usage = func() {
		flag.PrintDefaults()
	}

	flag.Parse()

	switch logFormat {
	case "text":
		env.Args.LogFormat = env.LogFormatText
	case "json":
		env.Args.LogFormat = env.LogFormatJSON
	default:
		slog.Error("Invalid log format", "format", logFormat)
		os.Exit(env.ExitCodeInvalidLogFormat)
	}
}

func initializeLogging() {
	var level slog.Leveler
	if env.Args.Debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	var handler slog.Handler
	if env.Args.LogFormat == "text" {
		handler = slog.NewTextHandler(
			os.Stderr, &slog.HandlerOptions{Level: level},
		)
	} else {
		handler = slog.NewJSONHandler(
			os.Stderr, &slog.HandlerOptions{Level: level},
		)
	}

	slog.SetDefault(slog.New(handler))
}

func startServer() {
	e := echo.New()

	// Initialize routes
	router(e)

	slog.Debug("Server started", "addr", env.Args.Addr, "server-path-prefix", env.Args.ServerPathPrefix)
	if err := e.Start(env.Args.Addr); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(env.ExitCodeServerStart)
	}
}
