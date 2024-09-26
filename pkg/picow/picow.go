package picow

const (
	GroupConfig = Group("config")
	GroupInfo   = Group("info")
	GroupLED    = Group("led")

	TypeSet = Type("set")
	TypeGet = Type("get")

	IDNoResponse = ID(-1)

	Port    = 3000
	EndByte = byte('\n')
)

var (
	Groups = []Group{
		GroupConfig,
		GroupInfo,
		GroupLED,
	}

	Types = []Type{
		TypeSet,
		TypeGet,
	}
)

type (
	Group string
	Type  string
	ID    int
)

type Request struct {
	Type    Type     `json:"type"`
	Group   Group    `json:"group"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	ID      ID       `json:"id"`
}

func NewRequest(id ID, t Type, g Group, c string, args ...string) *Request {
	if args == nil {
		args = make([]string, 0)
	}

	return &Request{
		ID:      id,
		Type:    t,
		Group:   g,
		Command: c,
		Args:    args,
	}
}

type Response struct {
	Data  any    `json:"data"`
	Error string `json:"error"`
	ID    ID     `json:"id"`
}

func NewResponse(id ID, data any, error string) *Response {
	return &Response{
		ID:    id,
		Data:  data,
		Error: error,
	}
}
