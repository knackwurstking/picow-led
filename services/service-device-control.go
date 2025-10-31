package services

import (
	"log/slog"

	"github.com/knackwurstking/picow-led/models"
)

// DeviceControl struct manages the control of devices.
type DeviceControl struct {
	registry *Registry
}

// NewDeviceControl creates a new instance of DeviceControl with the provided registry.
func NewDeviceControl(registry *Registry) *DeviceControl {
	return &DeviceControl{
		registry: registry,
	}
}

// CreateTable creates a table for device control data if it doesn't already exist.
// It returns an error if the database execution fails.
func (p *DeviceControl) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS device_control (
		device_id INTEGER PRIMARY KEY NOT NULL,
		color TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := p.registry.db.Exec(query)
	return err
}

func (p *DeviceControl) Get(deviceID models.DeviceID) (*models.DeviceControl, error) {
	slog.Debug("Get device control from database", "table", "device_control", "id", deviceID)

	query := `SELECT * FROM device_control WHERE device_id = ?`
	row := p.registry.db.QueryRow(query, deviceID)
	deviceControl, err := ScanDeviceControl(row)
	return deviceControl, HandleSqlError(err)
}

// Delete removes the record associated with a given device ID from the device_control table.
// It returns an error if the database execution fails or if the device is not found.
func (p *DeviceControl) Delete(deviceID models.DeviceID) error {
	slog.Debug("Delete device control from database", "table", "device_control", "id", deviceID)

	query := `DELETE FROM device_control WHERE device_id = ?`
	_, err := p.registry.db.Exec(query, deviceID)
	return err
}

// TODO: Methods for read and set current color, always store color not 0 in table. Get current color from the picow device directly

// TODO: Also handle version, temp and disk-usage

func ScanDeviceControl(scanner Scannable) (*models.DeviceControl, error) {
	control := &models.DeviceControl{}
	err := scanner.Scan(&control.DeviceID, &control.Color, &control.CreatedAt)
	return control, err
}
