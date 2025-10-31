package control

import "fmt"

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

func NewGetPinsRequest(id RequestID) *Request {
	return NewRequest(id, TypeGet, "config", "pins")
}

func NewSetPinsRequest(id RequestID, pins []uint8) *Request {
	var args []string
	for _, pin := range pins {
		args = append(args, fmt.Sprintf("%d", pin))
	}
	return NewRequest(
		id,
		TypeSet,
		"config",
		"pins",
		args...,
	)
}

func NewGetColorRequest(id RequestID) *Request {
	return NewRequest(id, TypeGet, "led", "color")
}

func NewSetColorRequest(id RequestID, color []uint8) *Request {
	var args []string
	for _, c := range color {
		args = append(args, fmt.Sprintf("%d", c))
	}
	return NewRequest(
		id,
		TypeSet,
		"led",
		"color",
		args...,
	)
}
