package routes

import (
	"log/slog"
	"picow-led/internal/api"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

func RegisterWS(e *echo.Echo, o *Options) {
	e.GET(o.ServerPathPrefix+"/ws", func(c echo.Context) error {
		websocket.Handler(func(conn *websocket.Conn) {
			defer conn.Close()

			client := &api.WSClient{
				Conn: conn,
			}

			o.WS.RegisterClient(client)
			defer o.WS.UnregisterClient(client)

			for {
				// Keep alive loop here, this websocket is readonly
				message := ""

				if err := websocket.Message.Receive(conn, &message); err != nil {
					slog.Warn("Receive websocket message failed",
						"error", err,
						"RealIP", c.RealIP(),
						"path", c.Request().URL.Path)
					break
				}

				slog.Debug("Receive websocket message",
					"message", message,
					"RealIP", c.RealIP(),
					"path", c.Request().URL.Path)
			}
		}).ServeHTTP(c.Response(), c.Request())

		return nil
	})
}
