package main

import "encoding/json"

type Request struct {
	Data    any     `json:"data"`
	Client  *client `json:"-"`
	Command string  `json:"command"`
}

func NewRequest(c *client, msg []byte) (*Request, error) {
	req := &Request{
		Client: c,
	}

	err := json.Unmarshal(msg, req)
	return req, err
}
