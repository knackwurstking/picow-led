package models

import "fmt"

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

func (c *Color) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(c.Duty) == 0 {
		return fmt.Errorf("duty cannot be empty")
	}
	return nil
}

var _ Model = (*Color)(nil)
