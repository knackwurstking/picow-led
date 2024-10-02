package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/MatusOllah/slogcolor"
	"github.com/SuperPaintman/nice/cli"
	"github.com/gorilla/websocket"
	"github.com/knackwurstking/picow-led-server/frontend"
	"github.com/labstack/echo/v4"
)

var (
	wsUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	connections = NewConnections()
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

				e := echo.New()

				// Init static file server
				e.StaticFS("/", frontend.GetFS())
				e.GET("/", func(c echo.Context) error {
					return c.Redirect(http.StatusSeeOther, "/index.html")
				})

				// Init (gorilla) websocket server
				e.GET("/api/ws", func(c echo.Context) error {
					conn, err := wsUpgrader.Upgrade(c.Response().Writer, c.Request(), nil)
					if err != nil {
						return c.String(http.StatusInternalServerError, err.Error())
					}

					// Add ws connection to connections struct (old clients struct)
					client := connections.Add(conn)
					defer connections.Delete(client.Conn)

					for {
						select {
						case data := <-client.Chan:
							client.Conn.SetWriteDeadline(client.WriteTimeout)
							if err := client.Conn.WriteJSON(data); err != nil {
								return c.String(http.StatusInternalServerError, err.Error())
							}
						case <-c.Request().Context().Done():
							return c.JSON(http.StatusOK, nil)
						case <-client.Done():
							return c.JSON(http.StatusOK, nil)
						}
					}
				})

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
