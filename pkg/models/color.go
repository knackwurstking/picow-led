package models

import (
	"fmt"
	"strings"
)

type Color struct {
	ID     ID       `json:"id"`
	Name   string   `json:"name"`
	Color  [3]uint8 `json:"color"`
	White  uint8    `json:"white,omitempty"`
	White2 uint8    `json:"white2,omitempty"`
}

func NewColor(name string, color [3]uint8, white, white2 uint8) *Color {
	return &Color{
		Name:   name,
		Color:  color,
		White:  white,
		White2: white2,
	}
}

func (c *Color) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}

func (c *Color) GetDuty(t DeviceType) []uint8 {
	duty := []uint8{}

	if t == DeviceTypeW {
		duty = append(duty, c.White)
	} else {
		duty = append(duty, c.Color[0], c.Color[1], c.Color[2])

		if t == DeviceTypeRGBW {
			duty = append(duty, c.White)
		}

		if t == DeviceTypeRGBWW {
			duty = append(duty, c.White2)
		}
	}

	return duty
}

func HexToColor(name, hex string) *Color {
	// Remove leading '#' if present
	hex = strings.TrimPrefix(hex, "#")

	var duty []uint8
	// Check if hex length is even
	if len(hex)%2 == 0 {
		for i := 0; i < len(hex); i += 2 {
			var d uint8
			fmt.Sscanf(hex[i:i+2], "%02x", &d)
			duty = append(duty, d)
		}
	} else {
		panic("invalid hex color: " + hex)
	}

	color := &Color{
		Name:   "",
		Color:  [3]uint8{0, 0, 0},
		White:  0,
		White2: 0,
	}

	if len(duty) == len(DeviceTypeW) {
		color.White = duty[0]
	} else if len(duty) >= len(DeviceTypeRGB) {
		color.Color[0] = duty[0]
		color.Color[1] = duty[1]
		color.Color[2] = duty[2]

		if len(duty) == len(DeviceTypeRGBW) {
			color.White = duty[3]
		}

		if len(duty) == len(DeviceTypeRGBWW) {
			color.White2 = duty[4]
		}
	}

	return color
}

func ColorToHex(duty [3]uint8) string {
	if len(duty) == 0 {
		return ""
	}

	hex := strings.Builder{}
	for _, d := range duty {
		fmt.Fprintf(&hex, "%02x", d)
	}

	return hex.String()
}

var _ Model = (*Color)(nil)
