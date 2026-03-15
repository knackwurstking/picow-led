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
	Color []uint8    `json:"color"` // TODO: Rename to Duty
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
			d.Color = []uint8{255, 255, 255} // Default to white if not provided
		}
	case DeviceTypeRGBW:
		if len(d.Color) != 4 {
			d.Color = []uint8{255, 255, 255, 255} // Default to white with no white component
		}
	case DeviceTypeRGBWW:
		if len(d.Color) != 5 {
			d.Color = []uint8{255, 255, 255, 255, 255} // Default to white with no white components
		}
	case DeviceTypeW:
		if len(d.Color) != 1 {
			d.Color = []uint8{255} // Default to full white if not provided
		}
	default:
		return fmt.Errorf("invalid device type: %s", d.Type)
	}
	return nil
}

func (d *Device) ToColor() *Color {
	color := &Color{}

	if d.Type == DeviceTypeW {
		color.White = d.Color[0]
		return color
	}

	if (d.Type == DeviceTypeRGB || d.Type == DeviceTypeRGBW || d.Type == DeviceTypeRGBWW) && len(d.Color) >= 3 {
		color.Color = [3]uint8{d.Color[0], d.Color[1], d.Color[2]}
	}

	switch d.Type {
	case DeviceTypeRGBW:
		if len(d.Color) >= 4 {
			color.White = d.Color[3]
		}
	case DeviceTypeRGBWW:
		if len(d.Color) >= 5 {
			color.White = d.Color[3]
			color.White2 = d.Color[4]
		}
	}

	return color
}

var _ Model = (*Device)(nil)
