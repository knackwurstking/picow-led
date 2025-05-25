package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"picow-led/internal/routes"
	"time"

	"github.com/SuperPaintman/nice/cli"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lmittmann/tint"
)

var (
	serverPathPrefix      = os.Getenv("SERVER_PATH_PREFIX")
	serverAddress         = os.Getenv("SERVER_ADDR")
	version               = "v0.1.0"
	apiConfigPath         = "api.yaml"
	apiConfigFallbackPath = ""
)

func init() {
	if d, err := os.UserConfigDir(); err == nil {
		apiConfigFallbackPath = filepath.Join(d, "picow-led", "api.yaml")
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

		// Register routes
		routesOptions := &routes.Options{
			ServerPathPrefix: serverPathPrefix,
		}

		routes.Register(e, routesOptions)

		return e.Start(*addr)
	}
}
