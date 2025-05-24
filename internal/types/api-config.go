package types

type APIConfig struct {
	Devices []*Device `json:"devices,omitempty" yaml:"devices,omitempty"`
}
