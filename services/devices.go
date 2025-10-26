package services

import (
	"log/slog"

	"github.com/knackwurstking/picow-led/models"
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

func (d *Devices) Get(id models.DeviceID) (*models.Device, error) {
	slog.Debug("Get device from database", "table", "devices", "id", id)

	query := `SELECT * FROM devices WHERE addr = ?`
	device, err := ScanDevice(d.registry.db.QueryRow(query, id))
	if err != nil {
		return nil, err
	}

	return device, nil
}

// TODO: ...
//func (d *Devices) List() ([]*models.Device, error)

// TODO: ...
//func (d *Devices) Add(device *models.Device) (models.DeviceID, error)

// TODO: ...
//func (d *Devices) Update(device *models.Device) error

// TODO: ...
//func (d *Devices) Delete(id models.DeviceID) error

func ScanDevice(r Scannable) (*models.Device, error) {
	var device models.Device
	err := r.Scan(&device.Addr, &device.Name, &device.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &device, nil
}
