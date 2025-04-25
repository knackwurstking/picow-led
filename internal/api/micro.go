package api

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	// Micro Default IDs

	MicroIDNoResponse = MicroID(-1)
	MicroIDDefault    = MicroID(0)

	// Micro Types

	MicroTypeGET = MicroType("get")
	MicroTypeSET = MicroType("set")

	// Micro Groups

	MicroGroupConfig = MicroGroup("config")
	MicroGroupLED    = MicroGroup("led")

	// MicroGroupInfo only used for `MicroTypeGET`
	MicroGroupInfo = MicroGroup("led")
)

type (
	MicroID    int
	MicroType  string
	MicroGroup string
)

type MicroRequest struct {
	MicroSocket

	ID    MicroID    `json:"id"`
	Group MicroGroup `json:"group"`
	Type  MicroType  `json:"type"`

	// Command mappings:
	//
	// Type: "set":
	// 		- Group: "config" Command: "led"
	// 		- Group: "led"    Command: "color"
	//
	// Type: "set":
	// 		- Group: "config" Command: "led"
	// 		- Group: "led"    Command: "color"
	// 		- Group: "info"   Command: "temp"
	// 		- Group: "info"   Command: "disk-usage"
	// 		- Group: "info"   Command: "version"
	Command string `json:"command"`

	// CommandArgs can be nil
	//
	// 	if type is `MicroTypeSET` and group is `MicroGroupConfig`
	//		if command is "led":
	// 			[]uint8 - range between 0-28
	//					  https://i0.wp.com/randomnerdtutorials.com/wp-content/uploads/2024/02/Raspberry-Pi-Pico-W-RP2040-Rev3-Board-Pinout-GPIOs.png?quality=100&strip=all&ssl=1
	//
	// 	elif type is `MicroTypeSET` and group is `MicroGroupLED`
	//		if command is "color":
	// 			[]uint8 - range between 0-255
	//
	// 	else
	// 		nil
	CommandArgs []string `json:"args"`
}

func (mr *MicroRequest) RequestPins(s *Server) (MicroPins, error) {
	if !mr.IsConnected() {
		if err := mr.Connect(s); err != nil {
			return nil, err
		}
	}

	// TODO: ...

	return nil, fmt.Errorf("under construction")
}

func (mr *MicroRequest) RequestColor(s *Server) (MicroColor, error) {
	// TODO: ...

	return nil, fmt.Errorf("under construction")
}

type (
	MicroPins      []uint
	MicroColor     []uint
	MicroTemp      int
	MicroDiskUsage struct {
		Used int `json:"used"`
		Free int `json:"free"`
	}
	MicroVersion string
)

type MicroSocket struct {
	socket net.Conn
	mutex  sync.Mutex
}

func (ms *MicroSocket) IsConnected() bool {
	return ms.socket != nil
}

func (ms *MicroSocket) Connect(s *Server) error {
	if ms.socket != nil {
		ms.Close()
	}

	dialer := net.Dialer{
		Timeout: time.Duration(time.Second * 5),
	}

	if conn, err := dialer.Dial("tcp", s.Addr); err != nil {
		ms.socket = nil
		s.Error = err.Error()
		s.Online = false
	} else {
		ms.socket = conn
		s.Error = ""
		s.Online = true
	}

	return errors.New(s.Error)
}

func (ms *MicroSocket) Read() ([]byte, error) {
	if ms.socket == nil {
		panic("socket is nil, call connect first")
	}

	data := make([]byte, 0)
	b := make([]byte, 1)
	endByte := []byte("\n")
	var n int
	var err error
	for {
		n, err = ms.socket.Read(b)
		if err != nil {
			break
		}
		if n == 0 {
			err = errors.New("no data")
			break
		}

		if bytes.Equal(b, endByte) {
			break
		}

		data = append(data, b...)
	}

	return data, err
}

func (ms *MicroSocket) Close() {
	if ms.socket == nil {
		return
	}

	ms.socket.Close()
	ms.socket = nil
}

type MicroResponse[T any | MicroPins | MicroColor | MicroTemp | MicroDiskUsage | MicroVersion] struct {
	ID    MicroID `json:"id"`
	Error string  `json:"error"`

	// Data contains the data requested
	Data T `json:"data"`
}
