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

// NewRequest Mappings:
//
//	`TypeSet` ("set"):
//		- group: "config" command: "pins"
//		- group: "led"    command: "color"
//
//	`TypeGet` ("get"):
//		- group: "config" command: "pins"
//		- group: "led"    command: "color"
//		- group: "info"   command: "temp"
//		- group: "info"   command: "disk-usage"
//		- group: "info"   command: "version"
//
// Args (optional):
//
//	if `Type` is "set" and `Group` is "config" and `Command` is "pins":
//		[]uint8 - range between 0-28 converted to a slice with strings
//			https://i0.wp.com/randomnerdtutorials.com/wp-content/uploads/2024/02/Raspberry-Pi-Pico-W-RP2040-Rev3-Board-Pinout-GPIOs.png?quality=100&strip=all&ssl=1
//
//	elif `Type` is "set" and `Group` is "led" and `Command` is "color":
//		[]uint8 - range between 0-255 converted to a slice with strings
//
//	else:
//		nil
//
// Examples:
//
//	`...(IDDefault, "get", "led", "color")`
//	`...(IDDefault, "set", "led", "color")`
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
