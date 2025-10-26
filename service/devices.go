package service

import (
	"log/slog"

	"github.com/knackwurstking/picow-led/model"
)

type Devices struct {
	registry *Registry
}

func NewDevices(r *Registry) *Devices {
	return &Devices{
		registry: r,
	}
}

func (d *Devices) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS devices (
		ID INTEGER PRIMARY KEY NOT NULL,
		addr TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := d.registry.db.Exec(query)
	return err
}

func (d *Devices) Get(id model.DeviceID) (*model.Device, error) {
	slog.Debug("Get device from database", "table", "devices", "id", id)

	query := `SELECT * FROM devices WHERE addr = ?`
	device, err := ScanDevice(d.registry.db.QueryRow(query, id))
	if err != nil {
		return nil, err
	}

	return device, nil
}

// TODO: ...
//func (d *Devices) List() ([]*model.Device, error)

// TODO: ...
//func (d *Devices) Add(device *model.Device) (model.DeviceID, error)

// TODO: ...
//func (d *Devices) Update(device *model.Device) error

// TODO: ...
//func (d *Devices) Delete(id model.DeviceID) error

func ScanDevice(r Scannable) (*model.Device, error) {
	var device model.Device
	err := r.Scan(&device.Addr, &device.Name, &device.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &device, nil
}
