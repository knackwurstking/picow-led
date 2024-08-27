package picow

const (
	GroupConfig = Group("config")
	GroupInfo   = Group("info")
	GroupLED    = Group("led")
	GroupMotion = Group("motion")

	TypeSet   = Type("set")
	TypeGet   = Type("get")
	TypeEvent = Type("event")

	IDNoResponse = ID(-1)

	DefaultPort    = 3000
	DefaultEndByte = byte('\n')
)

var (
	Groups = []Group{
		GroupConfig,
		GroupInfo,
		GroupLED,
		GroupMotion,
	}

	Types = []Type{
		TypeSet,
		TypeGet,
		TypeEvent,
	}
	Events = []string{
		"motion",
	}
)

// Group of command
type Group string

// Type of command
type Type string

// ID of command
type ID int

// Request object for the picow device
type Request struct {
	Group   Group    `json:"group"`
	Type    Type     `json:"type"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	ID      int      `json:"id"`
}

// Response object the picow device will respond with
type Response struct {
	Data  any    `json:"data"`
	Error string `json:"error"`
	ID    int    `json:"id"`
}
