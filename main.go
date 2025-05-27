package main

import (
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

var (
	ServerPathPrefix      = os.Getenv("SERVER_PATH_PREFIX")
	ServerAddress         = os.Getenv("SERVER_ADDR")
	Version               = "v0.1.0"
	ApiConfigPath         = "api.yaml"
	ApiConfigFallbackPath = ""
	DBPath                = "./database.db" // TODO: Use a default system path for this (not the config directory)
)

type Config struct {
	// TODO: Yaml configuration struct
}

func init() {
	if d, err := os.UserConfigDir(); err == nil {
		ApiConfigFallbackPath = filepath.Join(d, "picow-led", "api.yaml")
	}

	// TODO: Get dir to data storage (used for the database)
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
						cmd, "db",
						cli.WithShort("d"),
						cli.Usage("Change database"),
					)

					*addr = ServerAddress
					*dbPath = DBPath

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

func cliAction_Server(addr *string, dbPath *string) cli.ActionRunner {
	return func(cmd *cli.Command) error {
		e := echo.New()

		// Logger setup
		logger := slog.New(tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
			AddSource:  true,
		}))
		slog.SetDefault(logger)

		// Middleware
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
			Output: os.Stderr,
		}))

		// TODO: Load configuration and set devices to the database
		loadConfig()

		// Load API Configuration
		// slog.Info("Load API configuration",
		//	"path", apiConfigPath,
		//	"fallbackPath", apiConfigFallbackPath)

		//apiConfig, err := api.GetConfig(
		//	apiConfigPath, apiConfigFallbackPath,
		//)
		//if err != nil {
		//	slog.Warn("Read API configuration failed!")
		//	slog.Warn(err.Error())
		//}

		// Create and init the database
		db := database.NewDB(*dbPath)
		defer db.Close()

		// Register routes
		routesOptions := &routes.Options{
			ServerPathPrefix: ServerPathPrefix,
			DB:               db,
		}

		routes.Register(e, routesOptions)

		return e.Start(*addr)
	}
}

func loadConfig() *Config {
	// TODO: Continue here...
}
