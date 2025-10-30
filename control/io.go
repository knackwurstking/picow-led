package control

import (
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
//	  	- group: "led"    command: "color"
//
//	`TypeGet` ("get"):
//	 	- group: "config" command: "pins"
//	  	- group: "led"    command: "color"
//	  	- group: "info"   command: "temp"
//	  	- group: "info"   command: "disk-usage"
//	  	- group: "info"   command: "version"
//
// Args (optional):
//
//	if `Type` is "set" and `Group` is "config" and `Command` is "pins":
//		[]uint8 - range between 0-28 converted to a slice with strings
//			https://i0.wp.com/randomnerdtutorials.com/wp-content/uploads/2024/02/Raspberry-Pi-Pico-W-RP2040-Rev3-Board-Pinout-GPIOs.png?quality=100&strip=all&ssl=1
//
//	elif `Type` is "set" and `Group` is `led and `Command` is "color":
//		[]uint8 - range between 0-255 converted to a slice with strings
//
//	else
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

type Device struct {
	*models.Device
	*models.DeviceSetup

	conn net.Conn
}

func NewDevice(device *models.Device, setup *models.DeviceSetup) *Device {
	return &Device{
		Device:      device,
		DeviceSetup: setup,
	}
}

func (d *Device) Write(request []byte) (n int, err error) {
	// TODO: Connect, Send, Return
	dialer := net.Dialer{
		Timeout: time.Duration(time.Second * 5),
	}

	d.conn, err = dialer.Dial("tcp", string(d.Addr))
	if err != nil {
		return 0, err
	}

	n, err = d.conn.Write(request)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (d *Device) Read(response []byte) (n int, err error) {
	// TODO: Implement the Read method: Read, Parse Response, Return
	n, err = d.conn.Read(response)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (d *Device) Close() error {
	return d.conn.Close()
}

// EndByte returns the data with the end byte appended, only if not already present, newline will be used as end byte here
func (d *Device) EndByte(data []byte) []byte {
	if len(data) == 0 || data[len(data)-1] != '\n' {
		return append(data, '\n')
	}
	return data
}

var _ io.Writer = &Device{}
var _ io.Reader = &Device{}
