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
	return nil
}

var _ Model = (*Group)(nil)
