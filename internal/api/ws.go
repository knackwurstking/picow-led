package api

import "sync"

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

	ws.clients = append(ws.clients, c)
}

type WSClient struct{}
