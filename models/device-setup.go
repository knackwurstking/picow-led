package models

import "time"

type Pin uint8
type Pins []Pin

type DeviceSetup struct {
	DeviceID  DeviceID  `json:"device_id"`
	Pins      Pins      `json:"pins"`
	CreatedAt time.Time `json:"created_at"`
}

func NewDeviceSetup(deviceID DeviceID, pins Pins) *DeviceSetup {
	return &DeviceSetup{
		DeviceID:  deviceID,
		Pins:      pins,
		CreatedAt: time.Now(),
	}
}

func (p *DeviceSetup) Validate() bool {
	return p.DeviceID != 0 && len(p.Pins) > 0
}
