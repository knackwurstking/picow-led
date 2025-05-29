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
	ErrorCache = 2
)

var (
	ServerPathPrefix = os.Getenv("SERVER_PATH_PREFIX")
	ServerAddress    = os.Getenv("SERVER_ADDR")
	Version          = "v0.1.0"
	ConfigDir        string
	CacheDir         string
	ApiConfigPath    = "api.yaml"
	DBPath           = "database.db"
)

func init() {
	if d, err := os.UserConfigDir(); err == nil {
		ConfigDir = filepath.Join(d, "picow-led")
	}

	if d, err := os.UserCacheDir(); err != nil {
		CacheDir = filepath.Join(d, "picow-led")
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

					dbPath := cli.String(
						cmd, "cache",
						cli.WithShort("c"),
						cli.Usage("Cache location, where the database will land"),
					)

					*addr = ServerAddress
					*dbPath = filepath.Join(CacheDir)

					return cliAction_Server(addr, dbPath)
				}),
			},
		},
		CommandFlags: []cli.CommandFlag{
			cli.HelpCommandFlag(),
			cli.VersionCommandFlag(Version),
		},
	}

	app.HandleError(app.Run())
}

func cliAction_Server(addr *string, cache *string) cli.ActionRunner {
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
		dbPath := getDataBasePath(*cache)

		// Create database
		slog.Info(fmt.Sprintf("Database location: %s", dbPath), "CacheDir", CacheDir)
		db := database.NewDB(dbPath)
		defer db.Close()

		loadApiConfig(db)
		registerRoutes(e, db)

		return e.Start(*addr)
	}
}

func getDataBasePath(cache string) string {
	var err error

	if cache, err = filepath.Abs(cache); err != nil {
		slog.Error("Get the absolute database path failed", "error", err)
		os.Exit(ErrorCache)
	}

	dbPath := filepath.Join(cache, DBPath)

	if err = os.MkdirAll(filepath.Dir(dbPath), 0700); err != nil {
		slog.Error("Database creation failed", "error", err)
		os.Exit(ErrorCache)
	}

	return dbPath
}

func loadApiConfig(db *database.DB) {
	if config, err := loadConfig(ApiConfigPath, filepath.Join(ConfigDir, ApiConfigPath)); err != nil {
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
				slog.Debug("Database devices", "device", d)
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
		ServerPathPrefix: ServerPathPrefix,
		DB:               db,
	}

	routes.Register(e, routesOptions)
}
