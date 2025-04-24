// TODO: Either scan the network or read from configuration
package api

type Options struct {
	// TODO: Configuration stuff here
	Servers []Server `json,yaml:"servers"`
}

// Server contains host and port in use from a Device
type Server struct {
	Addr string `json,yaml:"addr"`
}

type Device struct {
	Server Server `json,yaml:"server"`
	// TODO: ...
}

func GetDevices(d Options) []Device {
	// TODO: ...

	return []Device{}
}
