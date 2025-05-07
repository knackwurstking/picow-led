package api

import (
	"encoding/json"
	"log"
	"slices"
	"sync"

	"golang.org/x/net/websocket"
)

type WS struct {
	clients []*WSClient

	broadcast chan any
	done      chan any

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

func (ws *WS) Start() error {
	for {
		select {
		case v := <-ws.broadcast:
			{
				wg := &sync.WaitGroup{}
				for _, c := range ws.clients {
					wg.Add(1)
					go func() {
						defer wg.Done()

						d, err := json.Marshal(v)
						if err != nil {
							log.Println(err, c.Conn)
						}

						err = websocket.Message.Send(c.Conn, d)
						if err != nil {
							log.Println(err, c.Conn)
						}
					}()
				}
				wg.Wait()
			}
		case <-ws.done:
			break
		}
	}
}

func (ws *WS) Broadcast(v any) {
	ws.broadcast <- v
}

func (ws *WS) Done() {
	ws.done <- nil
}

type WSClient struct {
	Conn *websocket.Conn
}
