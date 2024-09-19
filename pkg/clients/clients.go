package clients

import (
	"log/slog"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	EventTypeDevice  = EventType("device")
	EventTypeDevices = EventType("devices")
)

type EventType string

type Clients struct {
	Connections []*Client

	mutex *sync.Mutex
}

func NewClients() *Clients {
	return &Clients{
		Connections: make([]*Client, 0),
		mutex:       &sync.Mutex{},
	}
}

func (clients *Clients) Add(eventType EventType, conn *websocket.Conn) *Client {
	defer clients.mutex.Unlock()
	clients.mutex.Lock()

	for _, client := range clients.Connections {
		if client.Conn == conn && client.EventType == eventType {
			return client
		}
	}

	client := NewClient(eventType, conn)
	clients.Connections = append(clients.Connections, client)
	slog.Debug(
		"Added a client",
		"client.Conn.RemoveAddr()", client.Conn.RemoteAddr(),
		"len(clients)", len(clients.Connections),
	)

	return client
}

func (clients *Clients) Remove(eventType EventType, conn *websocket.Conn) {
	defer clients.mutex.Unlock()
	clients.mutex.Lock()

	for _, client := range clients.Connections {
		if client.Conn == conn && client.EventType == eventType {
			clients.removeClient(client)
			return
		}
	}
}

func (clients *Clients) Emit(eventType EventType, data any) {
	slog.Debug("Emit a new event", "eventType", eventType)

	for _, client := range clients.Connections {
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

	newConnections := make([]*Client, 0)
	for _, c := range clients.Connections {
		if c.Conn != client.Conn || c.EventType != client.EventType {
			newConnections = append(newConnections, c)
		}
	}
	clients.Connections = newConnections

	slog.Debug(
		"Removed a client",
		"client.Conn.RemoteAddr()", client.Conn.RemoteAddr(),
		"len(clients.Connections)", len(clients.Connections),
	)
}
