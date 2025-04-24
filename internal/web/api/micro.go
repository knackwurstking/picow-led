package api

type MicroRequest struct {
	ID      int    `json:"id"`
	Group   string `json:"group"`
	Type    string `json:"type"`
	Command string `json:"command"`

	// CommandArgs can be nil
	// TODO: Maybe using generics here, just like before
	CommandArgs []any `json:"args"`
}

type MicroResponse struct {
	ID    int    `json:"id"`
	Error string `json:"error"`

	// Data contains the data requested, can be nil
	// TODO: Maybe using generics here too, just like before
	Data any `json:"data"`
}
