package services

import (
	"database/sql"
	"fmt"
)

type Registry struct {
	db *sql.DB

	Device *DeviceService
	Color  *ColorService
	Group  *GroupService
}

func NewRegistry(db *sql.DB) (*Registry, error) {
	r := &Registry{db: db}

	r.Device = NewDeviceService(r)
	r.Color = NewColorService(r)
	r.Group = NewGroupService(r)

	if err := r.CreateTables(); err != nil {
		return nil, fmt.Errorf("create tables: %w", err)
	}

	return r, nil
}

func (r *Registry) CreateTables() error {
	services := []Service{
		r.Device,
		r.Color,
		r.Group,
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
