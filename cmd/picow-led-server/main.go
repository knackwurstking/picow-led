package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/MatusOllah/slogcolor"
	"github.com/SuperPaintman/nice/cli"
	"github.com/knackwurstking/picow-led-server/frontend"
	"github.com/knackwurstking/picow-led-server/internal/ws"
	"github.com/knackwurstking/picow-led-server/pkg/picow"
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

			var config *string
			cli.StringVar(cmd, config, "config",
				cli.Usage("Load api data from local json file"),
				cli.WithShort("c"),
				cli.Optional,
			)

			return func(cmd *cli.Command) error {
				// Initialize logger
				initLogger(debug, host, port)

				// Initialize api
				api := picow.NewApi()

				configPath, _ := os.UserConfigDir()
				path := filepath.Join(
					configPath, "picow-led-server", "api.json",
				)

				if err := api.LoadFromPath(path); err != nil {
					// Fallback path
					if err2 := api.LoadFromPath("api.json"); err2 != nil {
						slog.Warn(
							"Loading api configuration failed", "error", err,
						)
					}
				}

				// Init static file server
				public := frontend.GetFS()
				http.Handle("/", http.FileServerFS(public))

				// Init websocket handler
				room := ws.NewRoom(api)
				http.Handle("/ws", room)

				go room.Run()

				addr := fmt.Sprintf("%s:%d", host, port)
				slog.Info("Started server", "address", addr)
				return http.ListenAndServe(addr, &serverHandler{})
			}
		}),
		CommandFlags: []cli.CommandFlag{
			cli.HelpCommandFlag(),
			cli.VersionCommandFlag("0.7.0.dev"),
		},
	}

	app.HandleError(app.Run())
}

func initLogger(debug bool, host string, port uint) {
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
}
