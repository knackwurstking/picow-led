package main

import (
	"log/slog"

	"github.com/gorilla/websocket"
)

type client struct {
	socket  *websocket.Conn
	receive chan []byte
	room    *room
}

func (c *client) read() {
	defer c.close()

	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			slog.Debug(
				"Error while reading a message from a client",
				"client.address", c.socket.RemoteAddr(),
				"error", err,
			)
			return
		}

		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.close()

	for msg := range c.receive {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}

func (c *client) close() {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	c.socket.Close()
	close(c.receive)
}
