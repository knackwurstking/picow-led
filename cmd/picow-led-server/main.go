package main

import (
	"log"

	"github.com/SuperPaintman/nice/cli"
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
				// TODO: Init logger
				// TODO: Init static file server
				// TODO: Init socket io server
				// TODO: Start all of this

				log.Printf("Flags: debug=%+v, host=%+v, port=%+v", debug, host, port)
				return nil
			}
		}),
		CommandFlags: []cli.CommandFlag{
			cli.HelpCommandFlag(),
			cli.VersionCommandFlag("0.7.0.dev"),
		},
	}

	app.HandleError(app.Run())
}
