package api

type Options struct {
	Servers []Server `json,yaml:"servers"`
}

// Server contains host and port in use from a Device
type Server struct {
	Addr string `json,yaml:"addr"`
}

type Device struct {
	Server Server `json,yaml:"server"`
}

func GetDevices(d Options) []Device {
	// TODO: ...

	return []Device{}
}
