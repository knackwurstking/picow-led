package api

import (
	"slices"
	"sync"

	"golang.org/x/net/websocket"
)

type WS struct {
	clients []*WSClient

	mutex *sync.Mutex
}

func NewWS() *WS {
	return &WS{
		mutex: &sync.Mutex{},
	}
}

func (ws *WS) RegisterClient(c *WSClient) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	if slices.Contains(ws.clients, c) {
		return
	}

	ws.clients = append(ws.clients, c)
}

func (ws *WS) UnregisterClient(c *WSClient) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	for i, client := range ws.clients {
		if client == c {
			ws.clients = slices.Delete(ws.clients, i, 0)
		}
	}
}

func (ws *WS) Start() {
	// TODO: Create a room and listen for api events, broatcast to all
	// 		 connected clients, for now just devices data (cache) changes
}

type WSClient struct {
	Conn *websocket.Conn
}
