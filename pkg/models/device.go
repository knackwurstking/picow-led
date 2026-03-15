package models

import "fmt"

type Device struct {
	ID    ID      `json:"id"`
	Addr  string  `json:"addr"`
	Name  string  `json:"name"`
	Color []uint8 `json:"color"`
}

func NewDevice(addr, name string, color ...uint8) *Device {
	return &Device{
		Addr:  addr,
		Name:  name,
		Color: color,
	}
}

func (d *Device) Validate() error {
	if d.Addr == "" {
		return fmt.Errorf("device address cannot be empty")
	}
	return nil
}

var _ Model = (*Device)(nil)
