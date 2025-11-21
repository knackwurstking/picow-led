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
	"github.com/knackwurstking/picow-led/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lmittmann/tint"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := run()
	if err != nil {
		slog.Error("Application", "error", err)
		os.Exit(1)
	}
}

func run() error {
	err := parseFlags()
	if err != nil {
		return utils.WrapError(err, "parse flags")
	}

	switch env.Args.Command {
	case env.CommandServer:
		err = initializeLogging()
		if err != nil {
			return utils.WrapError(err, "initialize logging")
		}

		r, err := initializeDatabase()
		if err != nil {
			return utils.WrapError(err, "initialize database")
		}

		err = initializeDevices(r)
		if err != nil {
			return utils.WrapError(err, "initialize devices")
		}

		err = startServer(r)
		if err != nil {
			return utils.WrapError(err, "start server")
		}
	}

	return nil
}

func parseFlags() error {
	var logFormat string = string(env.Args.LogFormat)

	subCmd := flag.NewFlagSet("server", flag.ContinueOnError)

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
			return utils.WrapError(err, "parse server flags")
		}

		switch logFormat {
		case "text":
			env.Args.LogFormat = env.LogFormatText
		case "json":
			env.Args.LogFormat = env.LogFormatJSON
		default:
			return utils.WrapError(nil, "invalid log format: %s", logFormat)
		}

		err = verifyDatabasePath()
		if err != nil {
			return utils.WrapError(err, "database path validation")
		}

		env.Args.Command = env.CommandServer
	}

	return nil
}

func verifyDatabasePath() error {
	// Verify the database path
	if env.Args.DatabasePath == "" {
		return utils.WrapError(nil, "database path is required")
	}

	return nil
}

func initializeLogging() error {
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
	return nil
}

func initializeDatabase() (*services.Registry, error) {
	slog.Info("Initializing database", "path", env.Args.DatabasePath)

	sqlPath := fmt.Sprintf("%s", env.Args.DatabasePath)
	db, err := sql.Open("sqlite3", sqlPath)
	if err != nil {
		return nil, utils.WrapError(err, "open database connection: %s", sqlPath)
	}

	// Configure connection pool to handle multiple connections
	db.SetMaxOpenConns(25)   // Allow up to 25 open connections
	db.SetMaxIdleConns(25)   // Allow up to 25 idle connections
	db.SetConnMaxLifetime(0) // No maximum lifetime

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		return nil, utils.WrapError(err, "ping database")
	}

	r, err := services.NewRegistry(db)
	if err != nil {
		return nil, utils.WrapError(err, "create tables")
	}

	return r, nil
}

func initializeDevices(r *services.Registry) error {
	slog.Info("Initializing devices from the database")

	devices, err := r.Devices.List()
	if err != nil {
		return utils.WrapError(err, "list devices")
	}

	if len(devices) == 0 {
		return nil
	}

	wg := &sync.WaitGroup{}
	for _, device := range devices {
		wg.Go(func() {
			pins, err := r.DeviceControls.GetPins(device.ID)
			if err != nil {
				slog.Error("get device pins",
					"device_id", device.ID,
					"device_name", device.Name, "device_addr", device.Addr,
					"error", err)
				return
			}

			// This will get the color from the picow device and auto update the database
			currentColor, err := r.DeviceControls.GetCurrentColor(device.ID)
			if err != nil {
				slog.Error("get device color",
					"device_id", device.ID, "device", device.Name, "error", err)
				return
			}

			slog.Debug("Got device control data",
				"id", device.ID, "name", device.Name, "addr", device.Addr,
				"pins", pins, "current_color", currentColor)
		})
	}
	wg.Wait()

	return nil
}

func startServer(r *services.Registry) error {
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
		return utils.WrapError(err, "start server: %s", env.Args.Addr)
	}

	return nil
}
