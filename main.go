package main

import (
	"context"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"picow-led/components"
	"picow-led/web/js"
	"picow-led/web/pwa"

	"github.com/SuperPaintman/nice/cli"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

var (
	serverPathPrefix = os.Getenv("SERVER_PATH_PREFIX")
	serverAddr       = os.Getenv("SERVER_ADDR")
	version          = "v1.0.0"

	// NOTE: No need for a "/" at the beginning
	assetsToCache = []string{
		"",
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
)

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

		e.GET(serverPathPrefix+"/*", echo.StaticDirectoryHandler(dist(), false))
		e.GET(serverPathPrefix+"/", echo.WrapHandler(
			templ.Handler(
				components.Base(components.Data{
					ServerPathPrefix: serverPathPrefix,
					Version:          version,
				}), // TODO: Add page here
			),
		))

		// Server all "pwa" templates
		pwa.Serve(e, pwaTemplateData)

		// Server all "js" scripts
		js.Serve(e, jsTemplateData)

		return e.Start(*addr)
	}
}

func cliGenerateAction(path *string) cli.ActionRunner {
	return func(cmd *cli.Command) error {
		// Generate all templ stuff to `*path+"index.html"`
		err := os.MkdirAll(*path, 0700)
		if err != nil {
			return err
		}

		file, err := os.Create(filepath.Join(*path, "index.html"))
		if err != nil {
			return err
		}

		indexData := components.Data{
			ServerPathPrefix: serverPathPrefix,
			Version:          version,
		}
		err = components.Base(
			indexData,
			// TODO: Add page here
		).Render(context.Background(), file)
		if err != nil {
			return err
		}

		// Generate all PWA files
		err = pwa.Generate(*path, pwaTemplateData)
		if err != nil {
			return err
		}

		// Generate all JavaScript files
		err = js.Generate(*path, jsTemplateData)

		return nil
	}
}
