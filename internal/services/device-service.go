package services

import (
	"fmt"

	"github.com/knackwurstking/picow-led/internal/models"
)

type DeviceService struct {
	registry *Registry
}

func NewDeviceService(r *Registry) *DeviceService {
	return &DeviceService{
		registry: r,
	}
}

func (d *DeviceService) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS devices (
		id INTEGER NOT NULL,
		addr TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,

		PRIMARY KEY ("id" AUTOINCREMENT)
	);`

	_, err := d.registry.db.Exec(query)
	if err != nil {
		return NewServiceError("create devices table", err)
	}

	return nil
}

func (d *DeviceService) Get(deviceID models.ID) (*models.Device, error) {
	query := `SELECT * FROM devices WHERE id = ?`
	device, err := ScanDevice(d.registry.db.QueryRow(query, deviceID))
	if err != nil {
		return nil, NewServiceError("get device by ID", HandleSqlError(err))
	}

	return device, nil
}

func (d *DeviceService) GetByAddr(addr string) (*models.Device, error) {
	query := `SELECT * FROM devices WHERE addr = ?`
	device, err := ScanDevice(d.registry.db.QueryRow(query, addr))
	if err != nil {
		return nil, NewServiceError("get device by address", HandleSqlError(err))
	}

	return device, nil
}

func (d *DeviceService) List() ([]*models.Device, error) {
	rows, err := d.registry.db.Query(`SELECT * FROM devices`)
	if err != nil {
		return nil, NewServiceError("list devices", HandleSqlError(err))
	}
	defer rows.Close()

	var devices []*models.Device
	for rows.Next() {
		device, err := ScanDevice(rows)
		if err != nil {
			return nil, NewServiceError("scan device from rows", HandleSqlError(err))
		}
		devices = append(devices, device)
	}

	if err := rows.Err(); err != nil {
		return nil, NewServiceError("iterate device rows", err)
	}

	return devices, nil
}

func (d *DeviceService) Add(device *models.Device) (models.ID, error) {
	if device.Validate() != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidDevice, "device validation failed")
	}

	query := `INSERT INTO devices (addr, name) VALUES (?, ?)`
	result, err := d.registry.db.Exec(query, device.Addr, device.Name)
	if err != nil {
		return 0, NewServiceError("add device", HandleSqlError(err))
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, NewServiceError("get last inserted device ID", HandleSqlError(err))
	}

	return models.ID(id), nil
}

func (d *DeviceService) Update(device *models.Device) error {
	if device.Validate() != nil {
		return fmt.Errorf("%w: %v", ErrInvalidDevice, "device validation failed")
	}

	query := `UPDATE devices SET addr = ?, name = ? WHERE id = ?`
	_, err := d.registry.db.Exec(query, device.Addr, device.Name, device.ID)
	if err != nil {
		return NewServiceError("update device", HandleSqlError(err))
	}

	return nil
}

func (d *DeviceService) Delete(deviceID models.ID) error {
	query := `DELETE FROM device_controls WHERE device_id = ?`
	if err := d.registry.DeviceControl.Delete(deviceID); err != nil {
		return NewServiceError("delete device controls", HandleSqlError(err))
	}

	query = `DELETE FROM devices WHERE id = ?`
	if _, err := d.registry.db.Exec(query, deviceID); err != nil {
		return NewServiceError("delete device", HandleSqlError(err))
	}

	return nil
}

func ScanDevice(r Scannable) (*models.Device, error) {
	device := &models.Device{}
	err := r.Scan(&device.ID, &device.Addr, &device.Name)
	return device, err
}

var _ Service = (*DeviceService)(nil)
