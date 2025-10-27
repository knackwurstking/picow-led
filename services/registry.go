package services

import (
	"database/sql"
)

type Registry struct {
	db *sql.DB

	Devices      *Devices
	DeviceSetups *DeviceSetups
}

func NewRegistry(db *sql.DB) *Registry {
	r := &Registry{db: db}

	r.Devices = NewDevices(r)
	r.DeviceSetups = NewDeviceSetups(r)

	return r
}

func (r *Registry) CreateTables() error {
	var err error

	if err = r.Devices.CreateTable(); err != nil {
		return err
	}

	if err = r.DeviceSetups.CreateTable(); err != nil {
		return err
	}

	return nil
}
