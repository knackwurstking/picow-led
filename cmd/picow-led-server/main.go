package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/MatusOllah/slogcolor"
	"github.com/SuperPaintman/nice/cli"
	"github.com/knackwurstking/picow-led-server/frontend"
	"github.com/labstack/echo/v4"
)

func main() {
	app := cli.App{
		Name:  "picow-led-server",
		Usage: cli.Usage("PicoW LED Server"),
		Action: cli.ActionFunc(func(cmd *cli.Command) cli.ActionRunner {
			debug := false
			cli.BoolVar(cmd, &debug, "debug",
				cli.Usage("Enable debug logs"),
				cli.WithShort("d"),
				cli.Optional,
			)

			host := "0.0.0.0"
			cli.StringVar(cmd, &host, "host",
				cli.Usage("Change the default server host"),
				cli.WithShort("H"),
				cli.Optional,
			)

			port := uint(50833)
			cli.UintVar(cmd, &port, "port",
				cli.Usage("Change the default server port"),
				cli.WithShort("p"),
				cli.Optional,
			)

			return func(cmd *cli.Command) error {
				// Init logger
				if debug {
					slogcolor.DefaultOptions.Level = slog.LevelDebug
				}

				slog.SetDefault(
					slog.New(
						slogcolor.NewHandler(
							os.Stderr, slogcolor.DefaultOptions,
						),
					),
				)

				slog.Debug("Flags", "debug", debug, "host", host, "port", port)

				e := echo.New()

				// Init static file server
				e.StaticFS("/", frontend.GetFS())
				e.GET("/", func(c echo.Context) error {
					return c.Redirect(http.StatusSeeOther, "/index.html")
				})

				// TODO: Init (gorilla) websocket server

				return e.Start(fmt.Sprintf("%s:%d", host, port))
			}
		}),
		CommandFlags: []cli.CommandFlag{
			cli.HelpCommandFlag(),
			cli.VersionCommandFlag("0.7.0.dev"),
		},
	}

	app.HandleError(app.Run())
}
