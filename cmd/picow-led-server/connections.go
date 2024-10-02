package main

import (
	"log/slog"
	"sync"

	"github.com/gorilla/websocket"
)

type Connections struct {
	mutex *sync.Mutex
	conns []*Client
}

func NewConnections() *Connections {
	return &Connections{
		mutex: &sync.Mutex{},
	}
}

func (ws *Connections) Add(conn *websocket.Conn) *Client {
	slog.Debug("Add a connection",
		"remoteAddr", conn.RemoteAddr(),
		"localAddr", conn.LocalAddr(),
	)

	client := NewClient(conn)
	ws.conns = append(ws.conns, client)
	return client
}

func (ws *Connections) Delete(conn *websocket.Conn) {
	slog.Debug("Delete a connection",
		"remoteAddr", conn.RemoteAddr(),
		"localAddr", conn.LocalAddr(),
	)

	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	newConns := make([]*Client, 0)
	for _, c := range ws.conns {
		if c.Conn == conn {
			c.Conn.Close()
			continue
		}

		newConns = append(newConns, c)
	}

	ws.conns = newConns
}
