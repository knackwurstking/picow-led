// API Routes:
//   - GET - "/ws"
package routes

import (
	"picow-led/internal/api"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type wsOptions struct {
	ServerPathPrefix string
	WS               *api.WS
}

func wsRoutes(e *echo.Echo, o wsOptions) {
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
					c.Logger().Errorf("%s: %s", c.RealIP(), err)
					break
				}
				c.Logger().Debugf("%s: message: %s", c.RealIP(), message)
			}
		}).ServeHTTP(c.Response(), c.Request())

		return nil
	})
}
