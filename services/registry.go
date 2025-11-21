package services

import (
	"database/sql"

	"github.com/knackwurstking/picow-led/utils"
)

type Registry struct {
	db *sql.DB

	Devices        *Devices
	Colors         *Colors
	Groups         *Groups
	DeviceControls *DeviceControls
}

// NewRegistry creates a new registry instance and will call the `CreateTables` method.
func NewRegistry(db *sql.DB) (*Registry, error) {
	r := &Registry{db: db}

	r.Devices = NewDevices(r)
	r.Colors = NewColors(r)
	r.Groups = NewGroups(r)
	r.DeviceControls = NewDeviceControls(r)

	if err := r.CreateTables(); err != nil {
		return nil, utils.WrapError(err, "create tables")
	}

	return r, nil
}

func (r *Registry) CreateTables() error {
	services := []Service{
		r.Devices,
		r.Colors,
		r.Groups,
		r.DeviceControls,
	}

	var err error
	for _, service := range services {
		if err = service.CreateTable(); err != nil {
			return err
		}
	}

	return nil
}

func (r *Registry) Close() error {
	return r.db.Close()
}
