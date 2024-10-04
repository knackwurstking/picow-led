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

				// Init websocket handler
				room := ws.NewRoom()
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

	slog.Debug("Flags", "debug", debug, "host", host, "port", port)
}

type serverHandler struct{}

func (*serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered", "r", r)
		}
	}()

	crw := &serverResponseWriter{
		ResponseWriter: w,
		Hijacker:       w.(http.Hijacker),
	}

	crw.Header().Set("Access-Control-Allow-Origin", "*")
	crw.Header().Set("Content-Type", "text/plain")
	http.DefaultServeMux.ServeHTTP(crw, r)

	log := slog.Warn

	if crw.status >= 200 && crw.status < 300 {
		if crw.Header().Get("Content-Type") == "" {
			crw.Header().Set("Content-Type", "application/json")
		}

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
}

type serverResponseWriter struct {
	http.ResponseWriter
	http.Hijacker
	status int
}

func (crw *serverResponseWriter) WriteHeader(statusCode int) {
	crw.status = statusCode
	crw.ResponseWriter.WriteHeader(statusCode)
}
