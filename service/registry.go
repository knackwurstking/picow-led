package service

import (
	"database/sql"
)

type Registry struct {
	db *sql.DB

	Devices *Devices
}

func NewRegistry(db *sql.DB) *Registry {
	r := &Registry{db: db}

	r.Devices = NewDevices(r)

	return r
}

func (r *Registry) CreateTables() error {
	err := r.Devices.CreateTable()
	if err != nil {
		return err
	}

	return nil
}
