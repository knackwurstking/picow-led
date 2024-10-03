package main

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type room struct {
	clients map[*client]bool
	join    chan *client
	leave   chan *client
	forward chan []byte
}

func newRoom() *room {
	return &room{
		clients: make(map[*client]bool),
		join:    make(chan *client),
		leave:   make(chan *client),
		forward: make(chan []byte),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true

			slog.Debug(
				"Add a new client to the websocket room",
				"client.address", client.socket.RemoteAddr(),
				"clients", len(r.clients),
			)
		case client := <-r.leave:
			delete(r.clients, client)
			client.close()

			slog.Debug(
				"Remove a client from the websocket room",
				"client.address", client.socket.RemoteAddr(),
				"clients", len(r.clients),
			)
		case msg := <-r.forward:
			slog.Debug("Parse and handle a new command",
				"msg", string(msg))

			// TODO: Parse message (ex.: "GET api.devices"), start handler
			resp := msg // placeholder

			// TODO: Send response to (all) client(s?)
			slog.Debug("Forward a message to all clients",
				"clients", len(r.clients),
			)

			for client := range r.clients {
				client.receive <- resp
			}
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		slog.Error("ServeHTTP", "error", err)
		return
	}

	client := &client{
		socket:  socket,
		receive: make(chan []byte, messageBufferSize),
		room:    r,
	}

	r.join <- client
	defer func() {
		r.leave <- client
	}()

	go client.write()
	client.read()
}
