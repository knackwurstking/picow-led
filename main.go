package main

import (
	"context"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"picow-led/components"
	"picow-led/internal/config"
	"picow-led/internal/routes"
	"picow-led/web/js"
	"picow-led/web/pwa"

	"github.com/SuperPaintman/nice/cli"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	serverPathPrefix = os.Getenv("SERVER_PATH_PREFIX")
	serverAddr       = os.Getenv("SERVER_ADDR")
	version          = "v1.0.0"

	// NOTE: No need for a "/" at the beginning
	assetsToCache = []string{
		"",
		"settings",
		"manifest.json",
		"static/css/styles.css",
		"static/icons/apple-touch-icon-180x180.png",
		"static/icons/favicon.ico",
		"static/icons/icon.png",
		"static/icons/maskable-icon-512x512.png",
		"static/icons/pwa-192x192.png",
		"static/icons/pwa-512x512.png",
		"static/icons/pwa-64x64.png",
		"static/js/main.js",
		"static/screenshots/328x626.png",
		"static/screenshots/626x338.png",
	}

	pwaTemplateData = pwa.TemplateData{
		ServerPathPrefix: serverPathPrefix,
		Version:          version,
		AssetsToCache:    assetsToCache,
	}
	jsTemplateData = js.TemplateData{
		ServerPathPrefix: serverPathPrefix,
	}

	apiConfigPath         = "api.yaml"
	apiConfigFallbackPath = ""
)

func init() {
	if d, err := os.UserConfigDir(); err == nil {
		apiConfigFallbackPath = filepath.Join(d, "picow-led", "api.yaml")
	}
}

//go:embed dist
var _dist embed.FS

func dist() fs.FS {
	fs, err := fs.Sub(_dist, "dist")
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
			{
				Name:  "generate",
				Usage: cli.Usage("Generate HTML"),
				Action: cli.ActionFunc(func(cmd *cli.Command) cli.ActionRunner {
					path := cli.StringArg(
						cmd, "path",
						cli.Usage("destination directory"),
						cli.Required,
					)

					return cliGenerateAction(path)
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

		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
			Output: os.Stderr,
		}))

		e.GET(serverPathPrefix+"/*", echo.StaticDirectoryHandler(dist(), false))

		// Server all "pwa" templates
		pwa.Serve(e, pwaTemplateData)

		// Server all "js" scripts
		js.Serve(e, jsTemplateData)

		// The devices page will be at "/" for now
		e.GET(serverPathPrefix+"/", echo.WrapHandler(
			templ.Handler(
				components.Base(
					components.Data{
						ServerPathPrefix: serverPathPrefix,
						Version:          version,
					},
					components.PageDevices(),
				),
			),
		))

		// Settings page
		e.GET(serverPathPrefix+"/settings", echo.WrapHandler(
			templ.Handler(
				components.Base(
					components.Data{
						ServerPathPrefix: serverPathPrefix,
						Version:          version,
					},
					components.PageSettings(),
				),
			),
		))

		// Api
		e.Logger.Infof(
			"Read API configuration from: %s, %s",
			apiConfigPath, apiConfigFallbackPath,
		)
		routes.Create(e, routes.Options{
			ServerPathPrefix: serverPathPrefix,
			Api: config.GetApiOptions(
				apiConfigPath, apiConfigFallbackPath,
			),
		})

		return e.Start(*addr)
	}
}

type generatePage struct {
	filePath string
	page     templ.Component
}

func cliGenerateAction(path *string) cli.ActionRunner {
	return func(cmd *cli.Command) error {
		// Generate templ pages
		baseData := components.Data{
			ServerPathPrefix: serverPathPrefix,
			Version:          version,
		}

		pages := []generatePage{
			{filePath: "index.html", page: components.PageDevices()},
			{
				filePath: filepath.Join("settings", "index.html"),
				page:     components.PageSettings(),
			},
		}

		for _, p := range pages {
			file, err := createFile(*path, p.filePath)
			if err != nil {
				return err
			}

			err = components.Base(
				baseData, p.page,
			).Render(context.Background(), file)
			if err != nil {
				return err
			}
		}

		// Generate all PWA files
		err := pwa.Generate(*path, pwaTemplateData)
		if err != nil {
			return err
		}

		// Generate all JavaScript files
		err = js.Generate(*path, jsTemplateData)

		return nil
	}
}

func createFile(path, filePath string) (*os.File, error) {
	// Generate all templ stuff to `*path+"index.html"`
	fp := filepath.Join(path, filePath)
	_ = os.MkdirAll(filepath.Dir(fp), 0700)
	file, err := os.Create(fp)
	if err != nil {
		return nil, err
	}

	return file, err
}
