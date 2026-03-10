package models

import "time"

type DeviceSetup struct {
	DeviceID  DeviceID  `json:"device_id"`
	Pins      []uint8   `json:"pins"`
	CreatedAt time.Time `json:"created_at"`
}

func NewDeviceSetup(deviceID DeviceID, pins []uint8) *DeviceSetup {
	return &DeviceSetup{
		DeviceID:  deviceID,
		Pins:      pins,
		CreatedAt: time.Now(),
	}
}

func (p *DeviceSetup) Validate() bool {
	return p.DeviceID != 0 && len(p.Pins) > 0
}

var _ ServiceModel = (*DeviceSetup)(nil)
