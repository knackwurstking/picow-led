package models

import (
	"time"
)

type Device struct {
	ID        DeviceID  `json:"id"`
	Addr      Addr      `json:"addr"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func NewDevice(addr Addr, name string) *Device {
	return &Device{
		Addr:      addr,
		Name:      name,
		CreatedAt: time.Now(),
	}
}

func (d *Device) Validate() bool {
	return d.ID >= 0 && d.Addr != ""
}
