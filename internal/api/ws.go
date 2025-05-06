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

func (ws *WS) Start() {
	// TODO: Create a room and listen for api events, broatcast to all
	// 		 connected clients, for now just devices data (cache) changes
}

type WSClient struct{}
