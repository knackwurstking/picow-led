package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/MatusOllah/slogcolor"
	"github.com/SuperPaintman/nice/cli"
	"github.com/knackwurstking/picow-led-server/frontend"
	"github.com/knackwurstking/picow-led-server/internal/ws"
	"github.com/knackwurstking/picow-led-server/pkg/picow"
)

type Flags struct {
	Config *string
	Host   string
	Port   uint
	Debug  bool
}

func main() {
	app := cli.App{
		Name:  "picow-led-server",
		Usage: cli.Usage("PicoW LED Server"),
		Action: cli.ActionFunc(func(cmd *cli.Command) cli.ActionRunner {
			flags := &Flags{
				Host: "0.0.0.0",
				Port: uint(50833),
			}

			cli.BoolVar(cmd, &flags.Debug, "debug",
				cli.Usage("Enable debug logs"),
				cli.WithShort("d"),
				cli.Optional,
			)

			cli.StringVar(cmd, &flags.Host, "host",
				cli.Usage("Change the default server host"),
				cli.WithShort("H"),
				cli.Optional,
			)

			cli.UintVar(cmd, &flags.Port, "port",
				cli.Usage("Change the default server port"),
				cli.WithShort("p"),
				cli.Optional,
			)

			cli.StringVar(cmd, flags.Config, "config",
				cli.Usage("Load api data from local json file"),
				cli.WithShort("c"),
				cli.Optional,
			)

			return runCommand(flags)
		}),
		CommandFlags: []cli.CommandFlag{
			cli.HelpCommandFlag(),
			cli.VersionCommandFlag("0.7.0.dev"),
		},
	}

	app.HandleError(app.Run())
}

func runCommand(flags *Flags) cli.ActionRunner {
	return func(cmd *cli.Command) error {
		// Initialize logger
		if flags.Debug {
			slogcolor.DefaultOptions.Level = slog.LevelDebug
		}

		slog.SetDefault(
			slog.New(
				slogcolor.NewHandler(
					os.Stderr, slogcolor.DefaultOptions,
				),
			),
		)

		// Initialize api
		api := picow.NewApi()

		if flags.Config != nil {
			if *flags.Config == "" {
				*flags.Config = "api.json"
			}

			if err := api.LoadFromPath(*flags.Config); err != nil {
				slog.Warn("Loading api configuration failed", "error", err)
			}
		}

		// Init static file server
		public := frontend.GetFS()
		http.Handle("/", http.FileServerFS(public))

		// Init websocket handler
		room := ws.NewRoom(api)
		http.Handle("/ws", room)

		go room.Run()

		addr := fmt.Sprintf("%s:%d", flags.Host, flags.Port)
		slog.Info("Started server", "address", addr)
		return http.ListenAndServe(addr, &serverHandler{})
	}
}
