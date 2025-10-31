package control

import (
	"errors"
)

const (
	MinDuty uint8 = 0
	MaxDuty uint8 = 255

	MinPin uint8 = 0
	MaxPin uint8 = 15
)

var (
	ErrNotConnected = errors.New("not connected")
	ErrNoData       = errors.New("no data")
)

// TODO: Define types for parsing picow (microcontroller) device responses,
// and creating a request (JSON)

// Pins

func GetPins() ([]uint8, error)

func SetPins(pins ...uint8) error

// Duty

func GetColor() ([]uint8, error)

func SetColor(duty ...uint8) error
