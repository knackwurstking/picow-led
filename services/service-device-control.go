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

// Get retrieves a device control record from the database by device ID.
// It returns a DeviceControl instance and an error if the record is not found or execution fails.
func (p *DeviceControl) Get(deviceID models.DeviceID) (*models.DeviceControl, error) {
	slog.Debug("Get device control from database", "table", "device_control", "id", deviceID)

	query := `SELECT * FROM device_control WHERE device_id = ?`
	row := p.registry.db.QueryRow(query, deviceID)
	deviceControl, err := ScanDeviceControl(row)
	return deviceControl, HandleSqlError(err)
}

// List retrieves all device control records from the database.
// It returns a slice of DeviceControl instances and an error if execution fails.
func (p *DeviceControl) List() ([]*models.DeviceControl, error) {
	slog.Debug("List device controls from database", "table", "device_control")

	query := `SELECT * FROM device_control`
	rows, err := p.registry.db.Query(query)
	if err != nil {
		return nil, HandleSqlError(err)
	}
	defer rows.Close()

	deviceControls := make([]*models.DeviceControl, 0)
	for rows.Next() {
		deviceControl, err := ScanDeviceControl(rows)
		if err != nil {
			return nil, HandleSqlError(err)
		}
		deviceControls = append(deviceControls, deviceControl)
	}

	return deviceControls, nil
}

// Add inserts a new device control record into the database.
// It returns the inserted device ID and an error if execution fails.
func (p *DeviceControl) Add(deviceControl *models.DeviceControl) (models.DeviceID, error) {
	slog.Debug("Add device control to database", "table", "device_control", "device_id", deviceControl.DeviceID)

	query := `INSERT INTO device_control (device_id, color) VALUES (?, ?)`
	result, err := p.registry.db.Exec(query, deviceControl.DeviceID, deviceControl.Color)
	if err != nil {
		return 0, HandleSqlError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, HandleSqlError(err)
	}

	return models.DeviceID(id), nil
}

// Update updates an existing device control record in the database.
// It returns an error if execution fails or the record is not found.
func (p *DeviceControl) Update(deviceControl *models.DeviceControl) error {
	slog.Debug("Update device control in database", "table", "device_control", "id", deviceControl.DeviceID)

	query := `UPDATE device_control SET color = ? WHERE device_id = ?`
	_, err := p.registry.db.Exec(query, deviceControl.Color, deviceControl.DeviceID)
	return HandleSqlError(err)
}

// Delete removes the record associated with a given device ID from the database.
// It returns an error if the execution fails or if the device is not found.
func (p *DeviceControl) Delete(deviceID models.DeviceID) error {
	slog.Debug("Delete device control from database", "table", "device_control", "id", deviceID)

	query := `DELETE FROM device_control WHERE device_id = ?`
	_, err := p.registry.db.Exec(query, deviceID)
	return HandleSqlError(err)
}

// TODO: Methods for read and set current color, always store color not 0 in table. Get current color from the picow device directly

// TODO: Also handle version, temp and disk-usage

func ScanDeviceControl(scanner Scannable) (*models.DeviceControl, error) {
	control := &models.DeviceControl{}
	err := scanner.Scan(&control.DeviceID, &control.Color, &control.CreatedAt)
	return control, err
}
