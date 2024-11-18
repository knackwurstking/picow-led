package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/MatusOllah/slogcolor"
	"github.com/SuperPaintman/nice/cli"
	"github.com/knackwurstking/picow-led-server/frontend"
	"github.com/knackwurstking/picow-led-server/internal/models/ws"
	"github.com/knackwurstking/picow-led-server/pkg/picow"
)

type Flags struct{}

var flags = struct {
	config string
	host   string
	port   uint
	debug  bool
}{
	config: "api.json",
	host:   "0.0.0.0",
	port:   uint(50833),
}

func main() {
	app := cli.App{
		Name:  "picow-led-server",
		Usage: cli.Usage("PicoW LED Server"),
		Action: cli.ActionFunc(func(cmd *cli.Command) cli.ActionRunner {
			cli.BoolVar(cmd, &flags.debug, "debug",
				cli.Usage("Enable debug logs"),
				cli.WithShort("d"),
				cli.Optional,
			)

			cli.StringVar(cmd, &flags.host, "host",
				cli.Usage("Change the default server host"),
				cli.WithShort("H"),
				cli.Optional,
			)

			cli.UintVar(cmd, &flags.port, "port",
				cli.Usage("Change the default server port"),
				cli.WithShort("p"),
				cli.Optional,
			)

			cli.StringVar(cmd, &flags.config, "config",
				cli.Usage("Load api data from local json file"),
				cli.WithShort("c"),
				cli.Optional,
			)

			return runCommand
		}),
		// NOTE: Only "zsh" seems to be supported, maybe i should fork this
		//       repo and try to add some completions for "fish" and "bash"?
		// Commands: []cli.Command{
		// 	cli.CompletionCommand(),
		// },
		CommandFlags: []cli.CommandFlag{
			cli.HelpCommandFlag(),
			cli.VersionCommandFlag("0.8.3"),
		},
	}

	app.HandleError(app.Run())
}

func runCommand(cmd *cli.Command) error {
	// Initialize logger
	if flags.debug {
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

	if flags.config != "" {
		slog.Debug("Try to load configuration", "path", flags.config)
		if err := api.LoadFromPath(flags.config); err != nil {
			slog.Warn("Loading api configuration failed", "error", err)
		}
	}

	// Init static file server
	public := frontend.GetFS()
	http.Handle("GET /", http.FileServerFS(public))

	// Init websocket handler
	room := ws.NewRoom(api)

	if flags.config != "" {
		room.OnApiChange = func(api *picow.Api) {
			api.SaveToPath(flags.config)
		}
	}

	http.Handle("GET /ws", room)

	go room.Run()

	addr := fmt.Sprintf("%s:%d", flags.host, flags.port)
	slog.Info("Started server", "address", addr)
	return http.ListenAndServe(addr, &serverHandler{})
}
