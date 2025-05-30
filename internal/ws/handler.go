package ws

import (
	"slices"
	"sync"
)

type Handler struct {
	clients []*Client
	mutex   *sync.Mutex

	started bool

	chDone chan any
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

	h.chDone = make(chan any)
	defer func() {
		close(h.chDone)
	}()

	// TODO: Broadcast handler (chan)
	for {
		select {
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
