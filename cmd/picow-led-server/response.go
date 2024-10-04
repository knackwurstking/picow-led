package main

const (
	ResponseTypeDevices = "devices"
	ResponseTypeDevice  = "device"
)

type Response struct {
	Client *client
	Type   ResponseType
	Data   []byte
}

type ResponseType string
