package models

import (
	"fmt"
	"strings"
)

type Color struct {
	ID   ID      `json:"id"`
	Name string  `json:"name"`
	Duty []uint8 `json:"duty"`
}

func NewColor(name string, duty ...uint8) *Color {
	return &Color{
		Name: name,
		Duty: duty,
	}
}

func NewColorFromHex(name, hex string) *Color {
	var duty []uint8
	if len(hex) == 6 {
		for i := 0; i < 6; i += 2 {
			var d uint8
			fmt.Sscanf(hex[i:i+2], "%02x", &d)
			duty = append(duty, d)
		}
	}
	return &Color{
		Name: name,
		Duty: duty,
	}
}

func (c *Color) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(c.Duty) == 0 {
		return fmt.Errorf("duty cannot be empty")
	}
	return nil
}

func (c *Color) DutyToHex() string {
	if len(c.Duty) == 0 {
		return ""
	}

	hex := strings.Builder{}
	for _, d := range c.Duty[0:3] {
		fmt.Fprintf(&hex, "%02x", d)
	}
	return hex.String()
}

var _ Model = (*Color)(nil)
