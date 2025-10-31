package services

import "database/sql"

type Registry struct {
	db *sql.DB

	Devices        *Devices
	DeviceSetups   *DeviceSetups
	Colors         *Colors
	Groups         *Groups
	DeviceControls *DeviceControls
}

// NewRegistry creates a new registry instance and will call the `CreateTables` method.
func NewRegistry(db *sql.DB) *Registry {
	r := &Registry{db: db}

	r.Devices = NewDevices(r)
	r.DeviceSetups = NewDeviceSetups(r)
	r.Colors = NewColors(r)
	r.Groups = NewGroups(r)
	r.DeviceControls = NewDeviceControls(r)

	if err := r.CreateTables(); err != nil {
		panic("failed to create tables: " + err.Error())
	}

	return r
}

func (r *Registry) CreateTables() error {
	services := []Service{
		r.Devices,
		r.DeviceSetups,
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
