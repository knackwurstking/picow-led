package types

// Server contains host and port in use from a Device
type Server struct {
	Addr string `json:"addr" yaml:"addr"`
	// Name could be empty (optional)
	Name string `json:"name" yaml:"name"`
}
