package routes

import (
	"log/slog"
	"picow-led/internal/database"
	"picow-led/internal/ws"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type Options struct {
	ServerPathPrefix string
	DB               *database.DB
}

func Register(e *echo.Echo, o *Options) {
	apiHandler := NewAPIHandler(o.DB)
	registerGroup(e.Group(o.ServerPathPrefix+"/api"), apiHandler)
	registerWebSocket(e)
}

func registerGroup(g *echo.Group, apiHandler *APIHandler) {
	g.GET("/devices", apiHandler.GetDevices)

	g.GET("/devices/:addr", apiHandler.GetDevice)

	g.GET("/devices/:addr/name", apiHandler.GetDeviceName)

	g.GET("/devices/:addr/active_color", apiHandler.GetDeviceActiveColor)

	g.GET("/devices/:addr/color", apiHandler.GetDeviceColor)
	g.POST("/devices/:addr/color", apiHandler.PostDeviceColor)

	g.GET("/devices/:addr/pins", apiHandler.GetDevicePins)

	g.GET("/devices/:addr/power", apiHandler.GetDevicePower)
	g.POST("/devices/:addr/power", apiHandler.PostDevicePower)

	g.GET("/colors", apiHandler.GetColors)
	g.POST("/colors", apiHandler.PostColors)
	g.PUT("/colors", apiHandler.PutColors)

	g.GET("/colors/:id", apiHandler.GetColorsID)
	g.POST("/colors/:id", apiHandler.PostColorsID)
	g.DELETE("/colors/:id", apiHandler.DeleteColorsID)
}

func registerWebSocket(e *echo.Echo) {
	wsHandler := ws.NewHandler()

	e.GET("/ws", func(c echo.Context) error {
		websocket.Handler(func(conn *websocket.Conn) {
			defer conn.Close()

			client := ws.NewClient(conn)

			wsHandler.Register(client)
			defer wsHandler.Unregister(client)

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
