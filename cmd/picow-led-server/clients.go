package main

import (
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	EventTypeDevice  = EventType("device")
	EventTypeDevices = EventType("devices")
)

var (
	clientsMutex = &sync.Mutex{}
)

type EventType string

type Clients []*Client

func NewClients() Clients {
	return make(Clients, 0)
}

func (clients *Clients) Add(eventType EventType, conn *websocket.Conn) *Client {
	defer clientsMutex.Unlock()
	clientsMutex.Lock()

	for _, client := range *clients {
		if client.Conn == conn && client.EventType == eventType {
			return client
		}
	}

	client := NewClient(eventType, conn)
	*clients = append(*clients, client)
	slog.Debug("Added a client", "client", client, "len(clients)", len(*clients))

	return client
}

func (clients *Clients) Remove(eventType EventType, conn *websocket.Conn) {
	defer clientsMutex.Unlock()
	clientsMutex.Lock()

	for _, client := range *clients {
		if client.Conn == conn && client.EventType == eventType {
			clients.removeClient(client)
			return
		}
	}
}

func (clients *Clients) Emit(eventType EventType, data any) {
	slog.Debug("Emit a new event", "eventType", eventType)

	for _, client := range *clients {
		if client.EventType != eventType {
			continue
		}

		go func(client *Client) {
			client.Chan <- data
		}(client)
	}
}

func (clients *Clients) removeClient(client *Client) {
	defer client.Conn.Close()

	newClients := make([]*Client, 0)
	for _, c := range *clients {
		if c.Conn != client.Conn || c.EventType != client.EventType {
			newClients = append(newClients, c)
		}
	}
	*clients = newClients

	slog.Debug("Removed a client", "client", client, "len(*clients)", len(*clients))
}

type Client struct {
	EventType EventType
	Conn      *websocket.Conn
	Chan      chan any

	done chan any
}

func NewClient(eventType EventType, conn *websocket.Conn) *Client {
	return &Client{
		EventType: eventType,
		Conn:      conn,
		Chan:      make(chan any, 2),
		done:      make(chan any, 2),
	}
}

func (c *Client) StartHeartBeat() (exit chan any) {
	exit = make(chan any, 2)
	go func() {
		for {
			select {
			case <-exit:
				return
			default:
				time.Sleep(time.Millisecond * 2500)

				c.Conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 1000))
				if err := c.Conn.WriteMessage(websocket.BinaryMessage, []byte("ping")); err != nil {
					c.done <- nil
					return
				}

				c.Conn.SetReadDeadline(time.Now().Add(time.Millisecond * 2500))
				if _, p, err := c.Conn.ReadMessage(); (err == nil && string(p) != "pong") || err != nil {
					c.done <- nil
					return
				}
			}
		}
	}()

	return exit
}

func (c *Client) Done() chan any {
	return c.done
}
