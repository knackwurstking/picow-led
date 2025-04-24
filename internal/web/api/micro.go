package api

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
	ID      MicroID    `json:"id"`
	Group   MicroGroup `json:"group"`
	Type    MicroType  `json:"type"`
	Command string     `json:"command"`

	// CommandArgs can be nil
	// TODO: Maybe using generics here, just like before
	CommandArgs []any `json:"args"`
}

type MicroResponse struct {
	ID    MicroID `json:"id"`
	Error string  `json:"error"`

	// Data contains the data requested, can be nil
	// TODO: Maybe using generics here too, just like before
	Data any `json:"data"`
}
