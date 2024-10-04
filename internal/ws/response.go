package ws

const (
	ResponseTypeDevices = "devices"
	ResponseTypeDevice  = "device"
	ResponseTypeError   = "error"
)

type ResponseType string

type Response struct {
	Data any          `json:"data"`
	Type ResponseType `json:"type"`
}

func (r *Response) SetError(err error) {
	r.Type = ResponseTypeError
	r.Data = err.Error()
}
