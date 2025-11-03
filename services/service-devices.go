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
		id INTEGER PRIMARY KEY NOT NULL,
		addr TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := d.registry.db.Exec(query)
	return err
}

func (d *Devices) Get(id models.DeviceID) (*models.Device, error) {
	slog.Debug("Get device from database", "table", "devices", "id", id)

	query := `SELECT * FROM devices WHERE id = ?`
	device, err := ScanDevice(d.registry.db.QueryRow(query, id))
	if err != nil {
		return nil, HandleSqlError(err)
	}

	return device, nil
}

func (d *Devices) GetByAddr(addr models.Addr) (*models.Device, error) {
	slog.Debug("Get device from database", "table", "devices", "addr", addr)

	query := `SELECT * FROM devices WHERE addr = ?`
	device, err := ScanDevice(d.registry.db.QueryRow(query, addr))
	if err != nil {
		return nil, HandleSqlError(err)
	}

	return device, nil
}

func (d *Devices) List() ([]*models.Device, error) {
	slog.Debug("List devices from database", "table", "devices")

	rows, err := d.registry.db.Query(`SELECT * FROM devices`)
	if err != nil {
		return nil, HandleSqlError(err)
	}
	defer rows.Close()

	var devices []*models.Device
	for rows.Next() {
		device, err := ScanDevice(rows)
		if err != nil {
			return nil, HandleSqlError(err)
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func (d *Devices) Add(device *models.Device) (models.DeviceID, error) {
	slog.Debug("Add device to database", "table", "devices", "device", device)

	if !device.Validate() {
		return 0, ErrInvalidDevice
	}

	query := `INSERT INTO devices (addr, name) VALUES (?, ?)`
	result, err := d.registry.db.Exec(query, device.Addr, device.Name)
	if err != nil {
		return 0, HandleSqlError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, HandleSqlError(err)
	}

	return models.DeviceID(id), nil
}

func (d *Devices) Update(device *models.Device) error {
	slog.Debug("Update device in database", "table", "devices", "device", device)

	if !device.Validate() {
		return ErrInvalidDevice
	}

	query := `UPDATE devices SET addr = ?, name = ? WHERE id = ?`
	_, err := d.registry.db.Exec(query, device.Addr, device.Name, device.ID)
	if err != nil {
		return HandleSqlError(err)
	}

	return nil
}

func (d *Devices) Delete(id models.DeviceID) error {
	slog.Debug("Delete device from database", "table", "devices", "id", id)

	var err error

	query := `DELETE FROM device_controls WHERE device_id = ?`
	if err = d.registry.DeviceControls.Delete(id); err != nil {
		return HandleSqlError(err)
	}

	query = `DELETE FROM devices WHERE id = ?`
	if _, err = d.registry.db.Exec(query, id); err != nil {
		return HandleSqlError(err)
	}

	return nil
}

func ScanDevice(r Scannable) (*models.Device, error) {
	device := &models.Device{}
	err := r.Scan(&device.ID, &device.Addr, &device.Name, &device.CreatedAt)
	return device, err
}
