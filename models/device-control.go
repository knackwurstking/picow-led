package models

import (
	"time"

	"github.com/knackwurstking/picow-led/env"
)

type DeviceControl struct {
	DeviceID  DeviceID
	Color     []uint8
	CreatedAt time.Time
}

func NewDeviceControl(color []uint8) *DeviceControl {
	return &DeviceControl{
		Color: color,
	}
}

func (p *DeviceControl) Validate() bool {
	for _, c := range p.Color {
		if c < env.MinDuty || c > env.MaxDuty {
			return false
		}
	}
	return len(p.Color) > 0
}

var _ ServiceModel = (*DeviceControl)(nil)
