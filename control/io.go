package control

import (
	"bytes"
	"io"
	"net"
	"time"

	"github.com/knackwurstking/picow-led/models"
)

const (
	RequestIDNoResponse RequestID = -1
	RequestIDDefault    RequestID = 0
	TypeGet             Type      = "get"
	TypeSet             Type      = "set"
)

type (
	RequestID int
	Type      string
)

type Request struct {
	ID          RequestID `json:"id"`
	Type        Type      `json:"type"`
	Group       string    `json:"group"`
	Command     string    `json:"command"`
	CommandArgs []string  `json:"args"`
}

// ```python
//
//	COMMANDS = {
//	   "config": {
//	       "set": {
//	           "pins": config_set_pins, # []uint8
//	       },
//	       "get": {
//	           "pins": config_get_pins, # => []uint8
//	       },
//	   },
//	   "info": {
//	       "get": {
//	           "temp": info_get_temp,
//	           "disk-usage": info_get_disk_usage,
//	           "version": info_get_version,
//	       },
//	   },
//	   "led": {
//	       "set": {
//	           "color": led_set_color, # []uint8
//	       },
//	       "get": {
//	           "color": led_get_color, # => []uint8
//	       },
//	   },
//	}
//
// ```
func NewRequest(id RequestID, t Type, group string, command string, args ...string) *Request {
	return &Request{
		ID:          id,
		Type:        t,
		Group:       group,
		Command:     command,
		CommandArgs: args,
	}
}

type Response[T any] struct {
	ID    RequestID `json:"id"`
	Error string    `json:"error"`
	Data  T         `json:"data"`
}

type PicoW struct {
	*models.Device
	*models.DeviceSetup

	Conn net.Conn
}

func NewPicoW(device *models.Device, setup *models.DeviceSetup) *PicoW {
	return &PicoW{
		Device:      device,
		DeviceSetup: setup,
	}
}

func (p *PicoW) Write(request []byte) (n int, err error) {
	p.Connect()

	n, err = p.Conn.Write(p.EndByte(request))
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (p *PicoW) Read(response []byte) (n int, err error) {
	if p.Conn == nil {
		return 0, ErrNotConnected
	}

	n, err = p.Conn.Read(response)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (p *PicoW) ReadAll() (data []byte, err error) {
	if p.Conn == nil {
		return nil, ErrNotConnected
	}

	buffer := bytes.NewBuffer(make([]byte, 0))
	chunk := make([]byte, 1)
	for {
		p.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, err := p.Read(chunk)
		if err != nil {
			return nil, err
		}
		if n == 0 {
			return nil, ErrNoData
		}

		if bytes.Contains(chunk, []byte{'\n'}) {
			break
		}

		buffer.Write(chunk)
	}

	return buffer.Bytes(), nil
}

func (p *PicoW) Connect() error {
	if p.Conn != nil {
		return nil
	}

	dialer := net.Dialer{
		Timeout: time.Duration(time.Second * 5),
	}

	var err error
	p.Conn, err = dialer.Dial("tcp", string(p.Addr))
	if err != nil {
		return err
	}

	return nil
}

func (p *PicoW) Close() error {
	if err := p.Conn.Close(); err != nil {
		return err
	}

	p.Conn = nil
	return nil
}

// EndByte returns the data with the end byte appended, only if not already present, newline will be used as end byte here
func (p *PicoW) EndByte(data []byte) []byte {
	if len(data) == 0 || data[len(data)-1] != '\n' {
		return append(data, '\n')
	}
	return data
}

var _ io.Writer = &PicoW{}
var _ io.Reader = &PicoW{}
