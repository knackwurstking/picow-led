package clients

import (
	"log/slog"
	"sync"

	"github.com/gorilla/websocket"
)

type EmitType string

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

func (clients *Clients) Add(emitType EmitType, conn *websocket.Conn) *Client {
	defer clients.mutex.Unlock()
	clients.mutex.Lock()

	for _, client := range clients.Connections {
		if client.Conn == conn && client.EmitType == emitType {
			return client
		}
	}

	client := NewClient(emitType, conn)
	clients.Connections = append(clients.Connections, client)
	slog.Debug(
		"Added a client", "address", client.Conn.RemoteAddr(),
		"connections", len(clients.Connections),
	)

	return client
}

func (clients *Clients) Remove(emitType EmitType, conn *websocket.Conn) {
	defer clients.mutex.Unlock()
	clients.mutex.Lock()

	for _, client := range clients.Connections {
		if client.Conn == conn && client.EmitType == emitType {
			clients.removeClient(client)
			return
		}
	}
}

func (clients *Clients) Emit(emitType EmitType, data any) {
	slog.Debug("Emit a new event", "event", emitType)

	for _, client := range clients.Connections {
		if client.EmitType != emitType {
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
		if c.Conn != client.Conn || c.EmitType != client.EmitType {
			newConnections = append(newConnections, c)
		}
	}
	clients.Connections = newConnections

	slog.Debug(
		"Removed a client", "address", client.Conn.RemoteAddr(),
		"connections", len(clients.Connections),
	)
}
