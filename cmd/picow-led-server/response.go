package main

import "encoding/json"

const (
	ResponseTypeDevices = "devices"
	ResponseTypeDevice  = "device"
)

type ResponseType string

type Response struct {
	Data   any          `json:"data"`
	Client *client      `json:"-"`
	Type   ResponseType `json:"type"`
}

func (r *Response) JSON() []byte {
	data, _ := json.Marshal(r)
	return data
}
