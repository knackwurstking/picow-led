package ws

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/knackwurstking/picow-led-server/pkg/picow"
)

const (
	socketBufferSize = 1024
	// messageBufferSize = 1024
)

var (
	api      = picow.NewApi()
	upgrader = &websocket.Upgrader{
		ReadBufferSize:  socketBufferSize,
		WriteBufferSize: socketBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Room struct {
	clients   map[*Client]bool
	Join      chan *Client
	Leave     chan *Client
	Handle    chan *Request
	Broadcast chan *Response
}

func NewRoom() *Room {
	return &Room{
		clients:   make(map[*Client]bool),
		Join:      make(chan *Client),
		Leave:     make(chan *Client),
		Handle:    make(chan *Request),
		Broadcast: make(chan *Response),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Join:
			r.clients[client] = true

			slog.Debug(
				"Add a new client to the websocket room",
				"client.address", client.Socket.RemoteAddr(),
				"clients", len(r.clients),
			)
		case client := <-r.Leave:
			delete(r.clients, client)
			client.Close()

			slog.Debug(
				"Remove a client from the websocket room",
				"client.address", client.Socket.RemoteAddr(),
				"clients", len(r.clients),
			)
		case req := <-r.Handle:
			switch req.Data {
			case "GET api.devices":
				go func(req *Request) {
					req.Client.Response <- &Response{
						Type: ResponseTypeDevices,
						Data: api.Devices,
					}
				}(req)
			}
		case resp := <-r.Broadcast:
			for c := range r.clients {
				go func(c *Client) {
					c.Response <- resp
				}(c)
			}
		}
	}
}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		slog.Error("ServeHTTP", "error", err)
		return
	}

	client := NewClient(socket, r)
	r.Join <- client
	defer func() {
		r.Leave <- client
	}()

	go client.Write()
	client.Read()
}
