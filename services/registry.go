package services

import (
	"database/sql"
)

type Registry struct {
	db *sql.DB

	Devices *Devices
	Pins    *Pins
}

func NewRegistry(db *sql.DB) *Registry {
	r := &Registry{db: db}

	r.Devices = NewDevices(r)
	r.Pins = NewPins(r)

	return r
}

func (r *Registry) CreateTables() error {
	var err error

	if err = r.Devices.CreateTable(); err != nil {
		return err
	}

	if err = r.Pins.CreateTable(); err != nil {
		return err
	}

	return nil
}
