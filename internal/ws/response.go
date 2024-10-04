package ws

const (
	ResponseTypeDevices = "devices"
	ResponseTypeDevice  = "device"
)

type ResponseType string

type Response struct {
	Data any          `json:"data"`
	Type ResponseType `json:"type"`
}
