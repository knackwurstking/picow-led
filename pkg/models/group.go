package models

import "fmt"

type Group struct {
	ID      ID     `json:"id"`
	Name    string `json:"name"`
	Devices []ID   `json:"devices"`
}

func NewGroup(name string, devices ...ID) *Group {
	return &Group{
		Name:    name,
		Devices: devices,
	}
}

func (g *Group) Validate() error {
	if g.Name == "" {
		return fmt.Errorf("group name cannot be empty")
	}

	// Check devices, empty group not allowed
	if len(g.Devices) == 0 {
		return fmt.Errorf("group must contain at least one device")
	}

	return nil
}

var _ Model = (*Group)(nil)
