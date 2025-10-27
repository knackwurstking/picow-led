package models

import "time"

type Pins struct {
	DeviceID  DeviceID  `json:"device_id"`
	Name      string    `json:"name"`
	Pins      []uint8   `json:"pins"`
	CreatedAt time.Time `json:"created_at"`
}

func NewPins(deviceID DeviceID, name string, pins []uint8) *Pins {
	return &Pins{
		DeviceID:  deviceID,
		Name:      name,
		Pins:      pins,
		CreatedAt: time.Now(),
	}
}

func (p *Pins) Validate() bool {
	return p.DeviceID != 0 && p.Name != "" && len(p.Pins) > 0
}
