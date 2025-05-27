package micro

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
)

type Response[T any] struct {
	ID    ID     `json:"id"`
	Error string `json:"error"`
	Data  T      `json:"data"`
}

func Send(h *Handler, c *Command) ([]byte, error) {
	if !h.IsConnected() {
		if err := h.Connect(); err != nil {
			return nil, err
		}

		// NOTE: All connections created here will be closed, all other
		// 		 connections will stay open
		defer h.Close()
	}

	if c.CommandArgs == nil {
		c.CommandArgs = []string{}
	}

	data, err := json.Marshal(c)
	if err != nil {
		panic(err.Error())
	}
	err = h.Write(data)
	if err != nil {
		return nil, err
	}

	if c.ID == IDNoResponse {
		return []byte{}, nil
	}

	data, err = h.Read()
	if err != nil {
		return nil, err
	}

	slog.Debug("Got some data", "addr", h.Addr,
		"command", fmt.Sprintf("(%d) %s %s %s %#v",
			c.ID, c.Type, c.Group, c.Command, c.CommandArgs,
		),
	)

	return data, nil
}

func ParseResponse[T any](data []byte) (T, error) {
	resp := &Response[T]{}

	err := json.Unmarshal(data, resp)
	if err != nil {
		return resp.Data, err
	}

	if resp.Error != "" {
		return resp.Data, errors.New(resp.Error)
	}

	return resp.Data, nil
}

func GetColor(addr string) ([]uint8, error) {
	pkg := NewCommand(IDDefault, TypeGet, "led", "color")
	handler := NewHandler(addr)

	data, err := Send(handler, pkg)
	if err != nil {
		return nil, err
	}

	color, err := ParseResponse[[]uint8](data)
	if err != nil {
		return nil, err
	}

	return color, nil
}

func SetColor(addr string, color []uint8) error {
	pkg := NewCommand(IDDefault, TypeSet, "led", "color")

	// TODO: Send request and wait for response

	return errors.New("under construction")
}
