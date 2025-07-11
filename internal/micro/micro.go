package micro

import (
	"encoding/json"
	"errors"
	"strconv"
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

func GetPins(addr string) ([]int, error) {
	pkg := NewCommand(IDDefault, TypeGet, "config", "pins")
	handler := NewHandler(addr)

	data, err := Send(handler, pkg)
	if err != nil {
		return nil, err
	}

	pins, err := ParseResponse[[]int](data)
	if err != nil {
		return nil, err
	}

	return pins, nil
}

func GetColor(addr string) ([]int, error) {
	pkg := NewCommand(IDDefault, TypeGet, "led", "color")
	handler := NewHandler(addr)

	data, err := Send(handler, pkg)
	if err != nil {
		return nil, err
	}

	color, err := ParseResponse[[]int](data)
	if err != nil {
		return nil, err
	}

	return color, nil
}

func SetPins(addr string, pins []int) error {
	args := []string{}
	for _, n := range pins {
		args = append(args, strconv.Itoa(int(n)))
	}

	pkg := NewCommand(IDDefault, TypeSet, "config", "pins", args...)
	handler := NewHandler(addr)

	data, err := Send(handler, pkg)
	if err != nil {
		return err
	}

	_, err = ParseResponse[any](data)
	if err != nil {
		return err
	}

	return nil
}

func SetColor(addr string, color []int) error {
	args := []string{}
	for _, n := range color {
		args = append(args, strconv.Itoa(int(n)))
	}

	pkg := NewCommand(IDDefault, TypeSet, "led", "color", args...)
	handler := NewHandler(addr)

	data, err := Send(handler, pkg)
	if err != nil {
		return err
	}

	_, err = ParseResponse[any](data)
	if err != nil {
		return err
	}

	return nil
}
