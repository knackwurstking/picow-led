package main

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"picow-led/components"
	"picow-led/internal/api"
	"picow-led/internal/config"
	"picow-led/internal/routes"

	"github.com/SuperPaintman/nice/cli"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	serverPathPrefix      = os.Getenv("SERVER_PATH_PREFIX")
	serverAddr            = os.Getenv("SERVER_ADDR")
	version               = "v1.0.0"
	apiConfigPath         = "api.yaml"
	apiConfigFallbackPath = ""
)

func init() {
	if d, err := os.UserConfigDir(); err == nil {
		apiConfigFallbackPath = filepath.Join(d, "picow-led", "api.yaml")
	}
}

//go:embed public
var _public embed.FS

func public() fs.FS {
	fs, err := fs.Sub(_public, "public")
	if err != nil {
		panic(err)
	}

	return fs
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
						cli.Required,
					)

					*addr = serverAddr

					return cliServerAction(addr)
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

func cliServerAction(addr *string) cli.ActionRunner {
	return func(cmd *cli.Command) error {
		e := echo.New()

		// Echo: Middleware
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
			Output: os.Stderr,
		}))

		// Api Configuration
		apiOptions, err := config.GetApiOptions(
			apiConfigPath, apiConfigFallbackPath,
		)
		if err != nil {
			e.Logger.Warnf("Read API configuration failed: %s", err.Error())
		}

		// Echo: Static File Server
		e.GET(serverPathPrefix+"/*", echo.StaticDirectoryHandler(public(), false))

		// Base Data (templ)
		baseData := &components.BaseData{
			ServerPathPrefix: serverPathPrefix,
			Version:          version,
		}

		// Echo: / - page-devices
		e.GET(serverPathPrefix+"/", echo.WrapHandler(
			templ.Handler(
				components.Base(
					baseData,
					components.PageDevices(
						&components.PageDevicesData{
							BaseData: baseData,
							Devices:  api.GetDevices(apiOptions),
						},
					),
				),
			),
		))

		// Echo: /settings - page-settings
		e.GET(serverPathPrefix+"/settings", echo.WrapHandler(
			templ.Handler(
				components.Base(
					baseData,
					components.PageSettings(),
				),
			),
		))

		// Echo: Api
		e.Logger.Infof(
			"Read API configuration from: %s, %s",
			apiConfigPath, apiConfigFallbackPath,
		)

		routes.Create(e, routes.Options{
			ServerPathPrefix: serverPathPrefix,
			Api:              apiOptions,
		})

		return e.Start(*addr)
	}
}
