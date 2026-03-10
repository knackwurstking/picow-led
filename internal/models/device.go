package models

import "fmt"

type Device struct {
	ID   ID     `json:"id"`
	Addr string `json:"addr"`
	Name string `json:"name"`
}

func NewDevice(addr, name string) *Device {
	return &Device{
		Addr: addr,
		Name: name,
	}
}

func (d *Device) Validate() error {
	if d.Addr == "" {
		return fmt.Errorf("device address cannot be empty")
	}
	return nil
}

var _ Model = (*Device)(nil)
