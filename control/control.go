package control

import "errors"

const (
	MinRange uint8 = 0
	MaxRange uint8 = 255
)

var (
	ErrNotConnected = errors.New("not connected")
	ErrNoData       = errors.New("no data")
)

// TODO: Define types for parsing picow (microcontroller) device responses,
// and creating a request (JSON)

// Pins

//func GetPins() (models.Pins, error)
//func SetPins(pins models.Pins) error

// Duty

//func GetDuty() (models.Duty, error)
//func SetDuty(duty models.Duty) error
