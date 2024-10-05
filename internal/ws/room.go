package ws

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/knackwurstking/picow-led-server/pkg/picow"
)

const (
	socketBufferSize = 1024
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Room struct {
	Api       *picow.Api
	clients   map[*Client]bool
	Join      chan *Client
	Leave     chan *Client
	Handle    chan *Request
	Broadcast chan *Response
}

func NewRoom(api *picow.Api) *Room {
	return &Room{
		Api:       api,
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
			switch req.Command {

			case CommandGetApiDevices:
				go r.getApiDevices(req)

			case CommandPostApiDevice:
				go r.postApiDevice(req)

			case CommandPutApiDevice:
				// TODO: ... go putApiDevice(req)

			case CommandDeleteApiDevice:
				// TODO: ... go deleteApiDevice(req)

			case CommandPostApiDeviceColor:
				go r.postApiDeviceColor(req)
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

func (r *Room) getApiDevices(req *Request) {
	req.Client.Response <- &Response{
		Type: ResponseTypeDevices,
		Data: r.Api.Devices,
	}
}

func (r *Room) postApiDevice(req *Request) {
	if req.Data == "" {
		return
	}

	resp := &Response{}
	deviceData := picow.DeviceData{}
	if err := json.Unmarshal([]byte(req.Data), &deviceData); err != nil {
		resp.SetError(err)
		req.Client.Response <- resp
		return
	}

	device := picow.NewDevice(deviceData)
	r.Api.Devices = append(r.Api.Devices, device)

	resp.Set(ResponseTypeDevices, r.Api.Devices)
	r.Broadcast <- resp
}

func (r *Room) postApiDeviceColor(req *Request) {
	if req.Data == "" {
		return
	}

	var data struct {
		Addr  string      `json:"addr"`
		Color picow.Color `json:"color"`
	}

	resp := &Response{}

	if err := json.Unmarshal([]byte(req.Data), &data); err != nil {
		resp.SetError(err)
		req.Client.Response <- resp
		return
	}

	for _, d := range r.Api.Devices {
		if d.Addr() != data.Addr {
			continue
		}

		if err := d.SetColor(data.Color); err != nil {
			resp.SetError(err)
		} else {
			resp.Set(ResponseTypeDevice, d)
		}
	}

	if resp.Type == "" {
		resp.SetError(fmt.Errorf("device %s not found", data.Addr))
	}
	r.Broadcast <- resp
}
