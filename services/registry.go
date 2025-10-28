package services

import (
	"database/sql"
)

type Registry struct {
	db *sql.DB

	Devices      *Devices
	DeviceSetups *DeviceSetups
	Colors       *Colors
	Groups       *Groups
}

// NewRegistry creates a new registry instance and will call the `CreateTables` method.
func NewRegistry(db *sql.DB) *Registry {
	r := &Registry{db: db}

	r.Devices = NewDevices(r)
	r.DeviceSetups = NewDeviceSetups(r)
	r.Colors = NewColors(r)
	r.Groups = NewGroups(r)

	if err := r.CreateTables(); err != nil {
		panic("failed to create tables: " + err.Error())
	}

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

	if err = r.Colors.CreateTable(); err != nil {
		return err
	}

	if err = r.Groups.CreateTable(); err != nil {
		return err
	}

	return nil
}

func (r *Registry) Close() error {
	return r.db.Close()
}
