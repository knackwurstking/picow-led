package models

import "fmt"

type DeviceType string

const (
	DeviceTypeRGB   DeviceType = "RGB"
	DeviceTypeRGBW  DeviceType = "RGBW"
	DeviceTypeRGBWW DeviceType = "RGBWW"
	DeviceTypeW     DeviceType = "W"
)

type Device struct {
	ID    ID         `json:"id"`
	Addr  string     `json:"addr"`
	Name  string     `json:"name"`
	Color []uint8    `json:"color"`
	Type  DeviceType `json:"type"`
}

func NewDevice(addr, name string, t DeviceType, color ...uint8) *Device {
	return &Device{
		Addr:  addr,
		Name:  name,
		Color: color,
		Type:  t,
	}
}

func (d *Device) Validate() error {
	if d.Addr == "" {
		return fmt.Errorf("device address cannot be empty")
	}

	// Check the device type
	switch d.Type {
	case DeviceTypeRGB:
		if len(d.Color) != 3 {
			return fmt.Errorf("device color must have 3 components for RGB type")
		}
	case DeviceTypeRGBW:
		if len(d.Color) != 4 {
			return fmt.Errorf("device color must have 4 components for RGBW type")
		}
	case DeviceTypeRGBWW:
		if len(d.Color) != 5 {
			return fmt.Errorf("device color must have 5 components for RGBWW type")
		}
	case DeviceTypeW:
		if len(d.Color) != 1 {
			return fmt.Errorf("device color must have 1 component for W type")
		}
	default:
		return fmt.Errorf("invalid device type: %s", d.Type)
	}
	return nil
}

var _ Model = (*Device)(nil)
