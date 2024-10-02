package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/MatusOllah/slogcolor"
	"github.com/SuperPaintman/nice/cli"
	"github.com/knackwurstking/picow-led-server/frontend"
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
				initLogger(debug, host, port)

				// Init static file server
				public := frontend.GetFS()
				http.Handle("/", http.FileServerFS(public))

				// TODO: Init SocketIO Server...
				// ioServer := socketio.NewServer(&engineio.Options{})

				// ...

				// ioServer.Serve()
				// defer ioServer.Close()

				// http.Handle("/socket.io", ioServer)

				return http.ListenAndServe(
					fmt.Sprintf("%s:%d", host, port),
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						// TODO: Recover from a panic here?

						crw := &customResponseWriter{
							ResponseWriter: w,
						}

						http.DefaultServeMux.ServeHTTP(crw, r)

						log := slog.Warn
						if crw.status >= 200 && crw.status < 300 {
							log = slog.Info
						} else if crw.status >= 500 || crw.status == 0 {
							log = slog.Error
						}

						log("Request",
							"status", crw.status,
							"addr", r.RemoteAddr,
							"method", r.Method,
							"url", r.URL,
						)
					}),
				)
			}
		}),
		CommandFlags: []cli.CommandFlag{
			cli.HelpCommandFlag(),
			cli.VersionCommandFlag("0.7.0.dev"),
		},
	}

	app.HandleError(app.Run())
}

type customResponseWriter struct {
	http.ResponseWriter
	status int
}

func (crw *customResponseWriter) WriteHeader(statusCode int) {
	crw.status = statusCode
	crw.ResponseWriter.WriteHeader(statusCode)
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

	slog.Debug("Flags", "debug", debug, "host", host, "port", port)
}
