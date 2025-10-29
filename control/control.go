package control

import "github.com/knackwurstking/picow-led/models"

// TODO: Define types for parsing picow (microcontroller) device responses,
// and creating a request (JSON)

// Pins

func GetPins() (models.Pins, error)
func SetPins(pins models.Pins) error

// Range

// NOTE: Please verify the microcontroller repository first, as this feature might not be configurable per remote due to being hardcoded in the microcontroller.
func GetRange() (min, max uint8, err error)
func SetRange(min, max uint8) error

// Duty

func GetDuty() (models.Duty, error)
func SetDuty(duty models.Duty) error
