package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"picow-led/internal/database"
	"picow-led/internal/routes"
	"time"

	"github.com/SuperPaintman/nice/cli"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lmittmann/tint"
)

const (
	exitCodeCache = 2
)

var (
	serverPathPrefix = os.Getenv("SERVER_PATH_PREFIX")
	serverAddress    = os.Getenv("SERVER_ADDR")
	version          = "v0.13.1"
	configDir        string
	apiConfigPath    = "api.yaml"
	dbPath           = "database.db"
)

func init() {
	if d, err := os.UserConfigDir(); err == nil {
		configDir = filepath.Join(d, "picow-led")
	}
}

func main() {
	app := cli.App{
		Name:  "picow-led",
		Usage: cli.Usage("PicoW LED server and client."),
		Commands: []cli.Command{
			{
				Name:  "server",
				Usage: cli.Usage("Start the Server"),
				Action: cli.ActionFunc(func(cmd *cli.Command) cli.ActionRunner {
					addr := cli.String(
						cmd, "addr",
						cli.WithShort("a"),
						cli.Usage("Set server address (<host>:<port>)"),
					)

					*addr = serverAddress

					return cliAction_Server(addr)
				}),
			},
		},
		CommandFlags: []cli.CommandFlag{
			cli.HelpCommandFlag(),
			cli.VersionCommandFlag(version),
		},
	}

	app.HandleError(app.Run())
}

func cliAction_Server(addr *string) cli.ActionRunner {
	return func(cmd *cli.Command) error {
		e := echo.New()

		// Setup global logger
		logger := slog.New(tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
			AddSource:  true,
		}))
		slog.SetDefault(logger)

		middlewareHandlers(e)

		// Handle database path
		databasePath := filepath.Join(configDir, dbPath)
		createDatabasePath(databasePath)

		// Create database
		db := database.NewDB(databasePath)
		defer db.Close()

		StaticFS(e)

		loadApiConfig(db)
		registerRoutes(e, db)

		slog.Info("Server running",
			"addr", *addr, "serverPathPrefix", serverPathPrefix)
		return e.Start(*addr)
	}
}

func createDatabasePath(path string) {
	slog.Info(fmt.Sprintf("Database location: %s", path))

	var err error
	if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		slog.Error("Database creation failed", "error", err)
		os.Exit(exitCodeCache)
	}
}

func loadApiConfig(db *database.DB) {
	if config, err := loadConfig(apiConfigPath, filepath.Join(configDir, apiConfigPath)); err != nil {
		slog.Warn("Configuration", "error", err)
	} else {
		// Parse config and update database
		slog.Debug("Loaded configuration, Try update database...")
		devices := config.GetDataBaseDevices()
		if err = db.Devices.Set(devices...); err != nil {
			slog.Error("Set devices to database", "error", err)
		}

		// Log
		devices, err = db.Devices.List()
		if err != nil {
			slog.Error("Get device from database failed", "error", err)
		} else {
			for _, d := range devices {
				slog.Debug("Database device", "device", d)
			}
		}
	}
}

func middlewareHandlers(e *echo.Echo) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
		Output: os.Stderr,
	}))
}

func registerRoutes(e *echo.Echo, db *database.DB) {
	routesOptions := &routes.Options{
		ServerPathPrefix: serverPathPrefix,
		DB:               db,
	}

	routes.Register(e, routesOptions)
}
