package api

import "fmt"

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

type MicroResponse[T any | MicroPins | MicroColor | MicroTemp | MicroDiskUsage | MicroVersion] struct {
	ID    MicroID `json:"id"`
	Error string  `json:"error"`

	// Data contains the data requested
	Data T `json:"data"`
}
