package main

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize = 1024
	// messageBufferSize = 1024
)

var (
	api      = NewApi()
	upgrader = &websocket.Upgrader{
		ReadBufferSize:  socketBufferSize,
		WriteBufferSize: socketBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type room struct {
	clients   map[*client]bool
	join      chan *client
	leave     chan *client
	handle    chan *Request
	broadcast chan *Response
}

func newRoom() *room {
	return &room{
		clients:   make(map[*client]bool),
		join:      make(chan *client),
		leave:     make(chan *client),
		handle:    make(chan *Request),
		broadcast: make(chan *Response),
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
		case req := <-r.handle:
			switch req.Data {
			case "GET api.devices":
				go func(req *Request) {
					req.Client.response <- &Response{
						Client: req.Client,
						Type:   ResponseTypeDevices,
						Data:   api.Devices,
					}
				}(req)
			}
		case resp := <-r.broadcast:
			for c := range r.clients {
				go func(c *client) {
					c.response <- resp
				}(c)
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
		socket:   socket,
		response: make(chan *Response),
		room:     r,
	}

	r.join <- client
	defer func() {
		r.leave <- client
	}()

	go client.write()
	client.read()
}
