package models

import (
	"fmt"
	"slices"
)

type DeviceType string

const (
	DeviceTypeRGB   DeviceType = "RGB"
	DeviceTypeRGBW  DeviceType = "RGBW"
	DeviceTypeRGBWW DeviceType = "RGBWW"
	DeviceTypeW     DeviceType = "W"
)

type Device struct {
	ID   ID         `json:"id"`
	Addr string     `json:"addr"`
	Name string     `json:"name"`
	Duty []uint8    `json:"duty"`
	Type DeviceType `json:"type"`
}

func NewDevice(addr, name string, t DeviceType, duty ...uint8) *Device {
	return &Device{
		Addr: addr,
		Name: name,
		Duty: duty,
		Type: t,
	}
}

func (d *Device) Validate() error {
	if d.Addr == "" {
		return fmt.Errorf("device address cannot be empty")
	}

	// Check the device type
	switch d.Type {
	case DeviceTypeRGB:
		if len(d.Duty) != 3 {
			d.Duty = []uint8{255, 255, 255} // Default to white if not provided
		}
	case DeviceTypeRGBW:
		if len(d.Duty) != 4 {
			d.Duty = []uint8{255, 255, 255, 255} // Default to white with no white component
		}
	case DeviceTypeRGBWW:
		if len(d.Duty) != 5 {
			d.Duty = []uint8{255, 255, 255, 255, 255} // Default to white with no white components
		}
	case DeviceTypeW:
		if len(d.Duty) != 1 {
			d.Duty = []uint8{255} // Default to full white if not provided
		}
	default:
		return fmt.Errorf("invalid device type: %s", d.Type)
	}

	// Check if duty values are not all zero
	if slices.Max(d.Duty) == 0 {
		return fmt.Errorf("device duty cannot be all zero or empty")
	}

	return nil
}

func (d *Device) ToColor() *Color {
	color := &Color{}

	if d.Type == DeviceTypeW {
		color.White = d.Duty[0]
		return color
	}

	if (d.Type == DeviceTypeRGB || d.Type == DeviceTypeRGBW || d.Type == DeviceTypeRGBWW) && len(d.Duty) >= 3 {
		color.Color = [3]uint8{d.Duty[0], d.Duty[1], d.Duty[2]}
	}

	switch d.Type {
	case DeviceTypeRGBW:
		if len(d.Duty) >= 4 {
			color.White = d.Duty[3]
		}
	case DeviceTypeRGBWW:
		if len(d.Duty) >= 5 {
			color.White = d.Duty[3]
			color.White2 = d.Duty[4]
		}
	}

	return color
}

var _ Model = (*Device)(nil)
