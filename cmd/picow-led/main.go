package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/knackwurstking/picow-led/env"
	"github.com/knackwurstking/picow-led/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lmittmann/tint"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	parseFlags()

	switch env.Args.Command {
	case env.CommandServer:
		initializeLogging()
		r := initializeDatabase()
		initializeDevices(r)
		startServer(r)
	}
}

func parseFlags() {
	var logFormat string = string(env.Args.LogFormat)

	subCmd := flag.NewFlagSet("server", flag.ExitOnError)

	// Server Address
	subCmd.StringVar(&env.Args.Addr, "addr", env.Args.Addr, "Server address")

	// Server Path Prefix
	subCmd.StringVar(&env.Args.ServerPathPrefix, "path-prefix", env.Args.ServerPathPrefix, "Server path prefix")

	// Debug Flag
	subCmd.BoolVar(&env.Args.Debug, "debug", env.Args.Debug, "Enable debug mode")

	// Log Format: "text", "json"
	subCmd.StringVar(&logFormat, "log-format", logFormat, "Log format: text, json")

	// Get the required database path
	subCmd.StringVar(&env.Args.DatabasePath, "database-path", env.Args.DatabasePath, "Database path")

	// Custom Usage Message for subcommand
	flag.Usage = func() {
		fmt.Println("Usage: <program> [flags]")
		fmt.Println("Commands:")
		fmt.Println("\tserver\t\tStart the server")
	}

	// Custom Usage Message for server subcommand
	subCmd.Usage = func() {
		fmt.Println("Usage: <program> server [server-flags]")
		fmt.Println("Server Flags:")
		subCmd.PrintDefaults()
	}

	flag.Parse()

	if len(os.Args) > 1 && os.Args[1] == "server" {
		err := subCmd.Parse(os.Args[2:])
		if err != nil {
			slog.Error("Failed to parse flags", "error", err)
			os.Exit(env.ExitCodeInvalidFlags)
		}

		switch logFormat {
		case "text":
			env.Args.LogFormat = env.LogFormatText
		case "json":
			env.Args.LogFormat = env.LogFormatJSON
		default:
			slog.Error("Invalid log format", "format", logFormat)
			os.Exit(env.ExitCodeInvalidLogFormat)
		}

		verifyDatabasePath()

		env.Args.Command = env.CommandServer
	}
}

func verifyDatabasePath() {
	// Verify the database path
	if env.Args.DatabasePath == "" {
		slog.Error("Database path is required")
		os.Exit(env.ExitCodeInvalidDatabasePath)
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
		handler = tint.NewHandler(os.Stderr, &tint.Options{
			AddSource:  true,
			Level:      level,
			TimeFormat: time.DateTime,
		})
	} else {
		handler = slog.NewJSONHandler(
			os.Stderr, &slog.HandlerOptions{
				AddSource: true,
				Level:     level,
			},
		)
	}

	slog.SetDefault(slog.New(handler))
}

func initializeDatabase() *services.Registry {
	slog.Info("Initializing database", "path", env.Args.DatabasePath)

	// TODO: Move this database open and setup statement to the `services.NewRegistry` function
	sqlPath := fmt.Sprintf("%s", env.Args.DatabasePath)
	db, err := sql.Open("sqlite3", sqlPath)
	if err != nil {
		slog.Error("Failed to open database connection", "error", err)
		os.Exit(env.ExitCodeDatabaseConnection)
	}

	// Configure connection pool to prevent resource exhaustion
	db.SetMaxOpenConns(1)    // SQLite works best with single writer
	db.SetMaxIdleConns(1)    // Keep one connection alive
	db.SetConnMaxLifetime(0) // No maximum lifetime

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		slog.Error("Failed to ping database", "error", err)
		os.Exit(env.ExitCodeDatabasePing)
	}

	r := services.NewRegistry(db)
	if err := r.CreateTables(); err != nil {
		slog.Error("Failed to create tables", "error", err)
		os.Exit(env.ExitCodeDatabaseTables)
	}

	return r
}

func initializeDevices(r *services.Registry) {
	slog.Info("Initializing devices from the database")

	devices, err := r.Devices.List()
	if err != nil {
		slog.Error("Failed to list devices", "error", err)
		os.Exit(env.ExitCodeSetupDevices)
	}

	if len(devices) == 0 {
		return
	}

	wg := &sync.WaitGroup{}
	for _, device := range devices {
		wg.Go(func() {
			pins, err := r.DeviceControls.GetPins(device.ID)
			if err != nil {
				slog.Error("Failed to get device pins",
					"device_id", device.ID,
					"device_name", device.Name, "device_addr", device.Addr,
					"error", err)
				return
			}

			// This will get the color from the picow device and auto update the database
			currentColor, err := r.DeviceControls.GetCurrentColor(device.ID)
			if err != nil {
				slog.Error("Failed to get device color",
					"device_id", device.ID, "device", device.Name, "error", err)
				return
			}

			slog.Debug("Got device control data",
				"id", device.ID, "name", device.Name, "addr", device.Addr,
				"pins", pins, "current_color", currentColor)
		})
	}
	wg.Wait()
}

func startServer(r *services.Registry) {
	slog.Info("Starting server", "addr", env.Args.Addr, "prefix", env.Args.ServerPathPrefix)

	e := echo.New()

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output:           os.Stderr,
		Format:           "${time_custom} ${method} ${status} ${uri} ${latency_human} ${remote_ip} ${error}\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	// Initialize routes
	router(e, r)

	if err := e.Start(env.Args.Addr); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(env.ExitCodeServerStart)
	}
}
