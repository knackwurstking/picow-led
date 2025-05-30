package ws

import (
	"slices"
	"sync"
)

type Handler struct {
	clients []*Client
	mutex   *sync.Mutex
}

func NewHandler() *Handler {
	return &Handler{
		mutex: &sync.Mutex{},
	}
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
