package api

import (
	"encoding/json"
	"slices"
	"sync"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

const (
	BroadcastTypeDevices BroadcastType = "devices"
	BroadcastTypeColors  BroadcastType = "color"
)

type BroadcastType string

type BroadcastData struct {
	Type BroadcastType
	Data any
}

type WS struct {
	Clients []*WSClient

	logger echo.Logger

	running   bool
	broadcast chan BroadcastData
	done      chan any

	mutex *sync.Mutex
}

func NewWS(logger echo.Logger) *WS {
	return &WS{
		Clients: make([]*WSClient, 0),
		logger:  logger,
		mutex:   &sync.Mutex{},
	}
}

func (ws *WS) RegisterClient(c *WSClient) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	if slices.Contains(ws.Clients, c) {
		return
	}

	ws.Clients = append(ws.Clients, c)
}

func (ws *WS) UnregisterClient(c *WSClient) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	for i, client := range ws.Clients {
		if client == c {
			ws.Clients = slices.Delete(ws.Clients, i, 0)
		}
	}
}

func (ws *WS) Start() error {
	ws.running = true
	defer func() {
		ws.running = false
	}()

	ws.broadcast = make(chan BroadcastData)
	ws.done = make(chan any)

	defer func() {
		close(ws.broadcast)
		close(ws.done)
	}()

	for {
		select {
		case v := <-ws.broadcast:
			{
				wg := &sync.WaitGroup{}
				for _, c := range ws.Clients {
					wg.Add(1)
					go func() {
						defer wg.Done()

						d, err := json.Marshal(v)
						if err != nil && ws.logger != nil {
							ws.logger.Error(err, c)
						}

						err = websocket.Message.Send(c.Conn, d)
						if err != nil && ws.logger != nil {
							ws.logger.Warn(err, c)
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

func (ws *WS) Stop() {
	if !ws.running {
		return
	}

	ws.done <- nil
}

func (ws *WS) Broadcast(t BroadcastType, v any) {
	if !ws.running {
		return
	}

	ws.broadcast <- BroadcastData{
		Type: t,
		Data: v,
	}
}

type WSClient struct {
	Conn *websocket.Conn
}
