package api

type Device struct {
	Server Server `json,yaml:"server"`
}

func GetDevices(d Options) []Device {
	// TODO: ...

	return []Device{}
}
