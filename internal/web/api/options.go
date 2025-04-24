package api

type Options struct {
	Servers []Server `json,yaml:"servers"`
}

// Server contains host and port in use from a Device
type Server struct {
	Addr string `json,yaml:"addr"`
}
