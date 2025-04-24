package picow

type Request struct {
	Type    string   `json:"type"`
	Group   string   `json:"group"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	ID      ID       `json:"id"`
}
