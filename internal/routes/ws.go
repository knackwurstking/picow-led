// API Routes:
//   - TODO: GET - "/ws"
package routes

import (
	"picow-led/internal/api"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type wsOptions struct {
	ServerPathPrefix string
}

func wsRoutes(e *echo.Echo, ws *api.WS, o wsOptions) {
	e.GET(o.ServerPathPrefix+"/ws", func(c echo.Context) error {
		ws.RegisterClient(&api.WSClient{}) // TODO: ...

		websocket.Handler(func(c *websocket.Conn) {
			defer c.Close()
			for {
				// TODO: Keep alive loop here, this websocket is readonly
			}
		}).ServeHTTP(c.Response(), c.Request())

		return nil
	})
}
