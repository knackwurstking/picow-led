package control

import "github.com/knackwurstking/picow-led/models"

// Pins

func GetPins() (models.Pins, error)
func SetPins(pins models.Pins) error

// Range

func GetRange() (min, max uint8, err error)
func SetRange(min, max uint8) error

// Duty

func GetDuty() (models.Duty, error)
func SetDuty(duty models.Duty) error
