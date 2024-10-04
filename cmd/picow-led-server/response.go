package main

import "encoding/json"

const (
	ResponseTypeDevices = "devices"
	ResponseTypeDevice  = "device"
)

type ResponseType string

type Response struct {
	Client *client      `json:"-"`
	Type   ResponseType `json:"type"`
	Data   []byte       `json:"data"`
}

func (r *Response) JSON() []byte {
	data, _ := json.Marshal(r)
	return data
}
