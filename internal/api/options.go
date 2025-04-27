package api

type Config struct {
	Servers []*Server `json,yaml:"servers,omitempty"`
}

// Server contains host and port in use from a Device
type Server struct {
	Addr string `json,yaml:"addr"`
	// Name could be empty (optional)
	Name string `json,yaml:"name"`
}
