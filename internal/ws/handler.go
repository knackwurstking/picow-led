package ws

import (
	"encoding/json"
	"log/slog"
	"picow-led/internal/database"
	"slices"
	"sync"
)

const (
	BroadcastTypeDevices BroadcastType = "devices"
	BroadcastTypeDevice  BroadcastType = "device"
	BroadcastTypeColors  BroadcastType = "colors"
)

type BroadcastType string

type BroadcastData struct {
	Type BroadcastType `json:"type"`
	Data any           `json:"data"`
}

type Handler struct {
	clients []*Client
	mutex   *sync.Mutex

	started bool

	chBroadcast chan BroadcastData
	chDone      chan any
}

func NewHandler() *Handler {
	return &Handler{
		clients: make([]*Client, 0),
		mutex:   &sync.Mutex{},
	}
}

func (h *Handler) Start() {
	h.started = true
	defer func() {
		h.started = false
	}()

	h.chBroadcast = make(chan BroadcastData)
	h.chDone = make(chan any)
	defer func() {
		close(h.chDone)
	}()

	for {
		select {
		case v := <-h.chBroadcast:
			h.handleBroadcast(v)
		case <-h.chDone:
			break
		}
	}
}

func (h *Handler) Stop() {
	if !h.started {
		return
	}

	h.chDone <- nil
}

func (h *Handler) Register(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if slices.Contains(h.clients, client) {
		return
	}

	h.clients = append(h.clients, client)
}

func (h *Handler) Unregister(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for i, c := range h.clients {
		if c == client {
			h.clients = slices.Delete(h.clients, i, i+1)
		}
	}
}

func (h *Handler) handleBroadcast(v BroadcastData) {
	data, err := json.Marshal(v)
	if err != nil {
		slog.Error("(WS) Broadcast: Marshal JSON", "error", err, "v", v)
	}

	wg := &sync.WaitGroup{}
	for _, client := range h.clients {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err = client.Send(data); err != nil {
				slog.Warn("(WS) Broadcast: Send message",
					"error", err, "client", client)
			}
		}()
	}
	wg.Wait()
}

func (h *Handler) Broadcast(t BroadcastType, v any) {
	if !h.started {
		return
	}

	h.chBroadcast <- BroadcastData{
		Type: t,
		Data: v,
	}
}

func (h *Handler) BroadcastDevices(v []*database.Device) {
	if !h.started {
		return
	}

	h.chBroadcast <- BroadcastData{
		Type: BroadcastTypeDevices,
		Data: v,
	}
}

func (h *Handler) BroadcastDevice(v *database.Device) {
	if !h.started {
		return
	}

	h.chBroadcast <- BroadcastData{
		Type: BroadcastTypeDevice,
		Data: v,
	}
}

func (h *Handler) BroadcastColors(v []database.Color) {
	if !h.started {
		return
	}

	h.chBroadcast <- BroadcastData{
		Type: BroadcastTypeColors,
		Data: v,
	}
}
