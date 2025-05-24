package types

type Device struct {
	Server *Server `json:"server" yaml:"server"`

	Online bool   `json:"online" yaml:"-"` //  Not used in configurations
	Error  string `json:"error" yaml:"-"`  //  Not used in configurations

	// Color can be nil
	Color MicroColor `json:"color" yaml:"-"`
	// Pins can be nil
	Pins MicroPins `json:"pins" yaml:"pins"`
}
