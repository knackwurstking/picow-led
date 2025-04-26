package api

type Options struct {
	Servers []*Server `json,yaml:"servers,omitempty"`
}

// Server contains host and port in use from a Device
type Server struct {
	Addr string `json,yaml:"addr"`
	// Name could be empty (optional)
	Name   string `json,yaml:"name"`
	Online bool   `json:"online" yaml:"-"` //  Not used in configurations
	Error  string `json:"error" yaml:"-"`  //  Not used in configurations
}
