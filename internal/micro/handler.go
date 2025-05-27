package micro

import (
	"bytes"
	"errors"
	"net"
	"time"
)

type Handler struct {
	socket net.Conn
	addr   string
}

func NewHandler(addr string) *Handler {
	return &Handler{
		addr: addr,
	}
}

func (h *Handler) Addr() string {
	return h.addr
}

func (h *Handler) IsConnected() bool {
	return h.socket != nil
}

func (h *Handler) Connect() error {
	if h.socket != nil {
		h.Close()
	}

	dialer := net.Dialer{
		Timeout: time.Duration(time.Second * 5),
	}

	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		h.socket = nil
	} else {
		h.socket = conn
	}

	return err
}

func (h *Handler) Write(data []byte) error {
	data = append(data, []byte("\n")...)
	n, err := h.socket.Write(data)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("write: no data")
	}

	return nil
}

func (h *Handler) Read() ([]byte, error) {
	if h.socket == nil {
		panic("socket is nil, call connect first")
	}

	data := make([]byte, 0)
	b := make([]byte, 1)
	endByte := []byte("\n")
	var n int
	var err error
	for {
		n, err = h.socket.Read(b)
		if err != nil {
			break
		}
		if n == 0 {
			err = errors.New("read: no data")
			break
		}

		if bytes.Equal(b, endByte) {
			break
		}

		data = append(data, b...)
	}

	return data, err
}

func (h *Handler) Close() {
	if h.socket == nil {
		return
	}

	h.socket.Close()
	h.socket = nil
}
