package control

import "fmt"

const (
	RequestIDNoResponse RequestID = -1    // Indicates that no response is expected from the device.
	RequestIDDefault    RequestID = 0     // Default request ID.
	TypeGet             Type      = "get" // Indicates a GET command.
	TypeSet             Type      = "set" // Indicates a SET command.
)

type (
	RequestID int
	Type      string
)

// Request represents a command sent to the Picow device.
type Request struct {
	ID          RequestID `json:"id"`      // ID of the request, used for matching responses.
	Type        Type      `json:"type"`    // Type of command (GET or SET).
	Group       string    `json:"group"`   // Group of commands to which this request belongs.
	Command     string    `json:"command"` // Specific command within the group.
	CommandArgs []string  `json:"args"`    // Arguments for the command.
}

// NewRequest creates a new Request instance with the specified parameters.
func NewRequest(id RequestID, t Type, group string, command string, args ...string) *Request {
	if args == nil {
		args = make([]string, 0)
	}

	return &Request{
		ID:          id,
		Type:        t,
		Group:       group,
		Command:     command,
		CommandArgs: args,
	}
}

// NewGetPinsRequest creates a new GET request for the "pins" command.
func NewGetPinsRequest(id RequestID) *Request {
	return NewRequest(id, TypeGet, "config", "pins")
}

// NewSetPinsRequest creates a new SET request for the "pins" command with the specified pins.
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

// NewGetColorRequest creates a new GET request for the "color" command.
func NewGetColorRequest(id RequestID) *Request {
	return NewRequest(id, TypeGet, "led", "color")
}

// NewSetColorRequest creates a new SET request for the "color" command with the specified color values.
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

// NewGetTemperatureRequest creates a new GET request for the "temp" command.
func NewGetTemperatureRequest(id RequestID) *Request {
	return NewRequest(id, TypeGet, "info", "temp")
}

// NewGetDiskUsageRequest creates a new GET request for the "disk-usage" command.
func NewGetDiskUsageRequest(id RequestID) *Request {
	return NewRequest(id, TypeGet, "info", "disk-usage")
}

// NewGetVersionRequest creates a new GET request for the "version" command.
func NewGetVersionRequest(id RequestID) *Request {
	return NewRequest(id, TypeGet, "info", "version")
}
