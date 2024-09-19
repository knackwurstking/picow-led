package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	_clients "github.com/knackwurstking/picow-led-server/pkg/clients"
)

func endpointsEvents() {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	e.GET("/events/devices", func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		client := clients.Add(_clients.EventTypeDevices, conn)

		defer func() {
			clients.Remove(_clients.EventTypeDevices, conn)
		}()

		exit := client.StartHeartBeat()
		defer func() {
			exit <- nil
		}()

		for {
			select {
			case d := <-client.Chan:
				conn.SetWriteDeadline(time.Now().Add(time.Second))
				if err := conn.WriteJSON(d); err != nil {
					return c.JSON(http.StatusInternalServerError, err.Error())
				}
			case <-c.Request().Context().Done():
				return c.JSON(http.StatusOK, nil)
			case <-client.Done():
				return c.JSON(http.StatusOK, nil)
			}
		}
	})

	e.GET("/events/device", func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
		if err != nil {
			slog.Warn(
				fmt.Sprintf(
					"Internal server error %d: %s",
					http.StatusInternalServerError, err.Error(),
				),
			)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		client := clients.Add(_clients.EventTypeDevice, conn)

		defer func() {
			clients.Remove(_clients.EventTypeDevice, conn)
		}()

		exit := client.StartHeartBeat()
		defer func() {
			exit <- nil
		}()

		for {
			select {
			case d := <-client.Chan:
				conn.SetWriteDeadline(time.Now().Add(time.Second))
				if err := conn.WriteJSON(d); err != nil {
					return c.JSON(http.StatusInternalServerError, err.Error())
				}
			case <-c.Request().Context().Done():
				return c.JSON(http.StatusOK, nil)
			case <-client.Done():
				return c.JSON(http.StatusOK, nil)
			}
		}
	})
}
