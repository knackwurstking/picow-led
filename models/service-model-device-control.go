package models

import (
	"time"

	"github.com/knackwurstking/picow-led/env"
)

type DeviceControl struct {
	DeviceID   DeviceID  `json:"device_id"`
	Color      []uint8   `json:"color"`
	ModifiedAt time.Time `json:"modified_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func NewDeviceControl(id DeviceID, color []uint8) *DeviceControl {
	return &DeviceControl{
		DeviceID:   id,
		Color:      color,
		ModifiedAt: time.Now(),
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
