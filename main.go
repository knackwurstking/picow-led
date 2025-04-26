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

	//assetsToCache = []string{
	//	"",
	//	"settings",
	//	"manifest.json",
	//	"static/css/styles.css",
	//	"static/icons/apple-touch-icon-180x180.png",
	//	"static/icons/favicon.ico",
	//	"static/icons/icon.png",
	//	"static/icons/maskable-icon-512x512.png",
	//	"static/icons/pwa-192x192.png",
	//	"static/icons/pwa-512x512.png",
	//	"static/icons/pwa-64x64.png",
	//	"static/js/main.js",
	//	"static/js/api.js",
	//	"static/js/page-devices.js",
	//	"static/screenshots/328x626.png",
	//	"static/screenshots/626x338.png",
	//}
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
			//{
			//	Name:  "generate",
			//	Usage: cli.Usage("Generate HTML"),
			//	Action: cli.ActionFunc(func(cmd *cli.Command) cli.ActionRunner {
			//		path := cli.StringArg(
			//			cmd, "path",
			//			cli.Usage("destination directory"),
			//			cli.Required,
			//		)

			//		return cliGenerateAction(path)
			//	}),
			//},
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

type generatePage struct {
	filePath string
	page     templ.Component
}

//func cliGenerateAction(path *string) cli.ActionRunner {
//	return func(cmd *cli.Command) error {
//		// Generate templ pages
//		baseData := &components.BaseData{
//			ServerPathPrefix: serverPathPrefix,
//			Version:          version,
//		}
//
//		pages := []generatePage{
//			{
//				filePath: "index.html",
//				page: components.PageDevices(
//					&components.PageDevicesData{
//						BaseData: baseData,
//					},
//				),
//			},
//			{
//				filePath: filepath.Join("settings", "index.html"),
//				page:     components.PageSettings(),
//			},
//		}
//
//		for _, p := range pages {
//			file, err := createFile(*path, p.filePath)
//			if err != nil {
//				return err
//			}
//
//			err = components.Base(
//				baseData, p.page,
//			).Render(context.Background(), file)
//			if err != nil {
//				return err
//			}
//		}
//
//		return nil
//	}
//}
//
//func createFile(path, filePath string) (*os.File, error) {
//	// Generate all templ stuff to `*path+"index.html"`
//	fp := filepath.Join(path, filePath)
//	_ = os.MkdirAll(filepath.Dir(fp), 0700)
//	file, err := os.Create(fp)
//	if err != nil {
//		return nil, err
//	}
//
//	return file, err
//}
