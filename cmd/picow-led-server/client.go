package main

import (
	"log/slog"

	"github.com/gorilla/websocket"
)

type client struct {
	socket   *websocket.Conn
	response chan *Response
	room     *room
}

func (c *client) read() {
	defer c.socket.Close()

	for {
		mt, msg, err := c.socket.ReadMessage()
		if err != nil {
			slog.Debug(
				"Error while reading a message from a client",
				"client.address", c.socket.RemoteAddr(),
				"error", err,
			)
			return
		}

		slog.Debug(
			"Got a message from a client",
			"client.address", c.socket.RemoteAddr(),
			"message.type", mt,
		)

		req, err := NewRequest(c, msg)
		if err != nil {
			slog.Warn("Parsing request failed", "error", err)
			continue
		}
		c.room.handle <- req
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for resp := range c.response {
		err := c.socket.WriteMessage(websocket.TextMessage, resp.JSON())
		if err != nil {
			return
		}
	}
}

func (c *client) close() {
	close(c.response)
}
