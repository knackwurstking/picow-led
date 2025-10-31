package services

import (
	"log/slog"

	"github.com/knackwurstking/picow-led/control"
	"github.com/knackwurstking/picow-led/models"
)

// DeviceSetups struct manages the setup of devices.
type DeviceSetups struct {
	registry *Registry
}

// NewDeviceSetups creates a new instance of DeviceSetups with the provided registry.
func NewDeviceSetups(registry *Registry) *DeviceSetups {
	return &DeviceSetups{
		registry: registry,
	}
}

// CreateTable creates a table for device setup data if it doesn't already exist.
// It returns an error if the database execution fails.
func (p *DeviceSetups) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS device_setups (
		device_id INTEGER UNIQUE NOT NULL,
		pins TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := p.registry.db.Exec(query)
	return err
}

// Get retrieves a device setup record from the database by device ID.
// It returns a DeviceSetup instance and an error if the record is not found or execution fails.
func (p *DeviceSetups) Get(id models.DeviceID) (*models.DeviceSetup, error) {
	slog.Debug("Get device setup from database", "table", "device_setups", "id", id)

	query := `SELECT * FROM device_setups WHERE device_id = ?`
	devicesSetup, err := ScanDevicesSetup(p.registry.db.QueryRow(query, id))
	if err != nil {
		return nil, HandleSqlError(err)
	}

	return devicesSetup, nil
}

// List retrieves all device setup records from the database.
// It returns a slice of DeviceSetup instances and an error if execution fails.
func (p *DeviceSetups) List() ([]*models.DeviceSetup, error) {
	slog.Debug("List device setups from database", "table", "device_setups")

	query := `SELECT * FROM device_setups`
	rows, err := p.registry.db.Query(query)
	if err != nil {
		return nil, HandleSqlError(err)
	}
	defer rows.Close()

	deviceSetups := make([]*models.DeviceSetup, 0)
	for rows.Next() {
		devicesSetup, err := ScanDevicesSetup(rows)
		if err != nil {
			return nil, HandleSqlError(err)
		}
		deviceSetups = append(deviceSetups, devicesSetup)
	}

	return deviceSetups, nil
}

// Add inserts a new device setup record into the database.
// It returns the inserted device ID and an error if execution fails.
func (p *DeviceSetups) Add(deviceSetup *models.DeviceSetup) (models.DeviceID, error) {
	slog.Debug("Add device setup to database", "table", "device_setups", "device_id", deviceSetup.DeviceID)

	if !deviceSetup.Validate() {
		return 0, ErrInvalidDeviceSetup
	}

	query := `INSERT INTO device_setups (device_id, pins) VALUES (?, ?)`
	result, err := p.registry.db.Exec(query, deviceSetup.DeviceID, deviceSetup.Pins)
	if err != nil {
		return 0, HandleSqlError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, HandleSqlError(err)
	}

	return models.DeviceID(id), p.controlPins(deviceSetup.DeviceID, deviceSetup.Pins)
}

// Update updates an existing device setup record in the database.
// It returns an error if execution fails or the record is not found.
func (p *DeviceSetups) Update(deviceSetup *models.DeviceSetup) error {
	slog.Debug("Update device setup in database", "table", "device_setups", "deviceSetup", deviceSetup)

	if !deviceSetup.Validate() {
		return ErrInvalidDeviceSetup
	}

	query := `UPDATE device_setups SET pins = ? WHERE device_id = ?`
	_, err := p.registry.db.Exec(query, deviceSetup.Pins, deviceSetup.DeviceID)
	if err != nil {
		return HandleSqlError(err)
	}

	return p.controlPins(deviceSetup.DeviceID, deviceSetup.Pins)
}

// Delete removes the record associated with a given device ID from the database.
// It returns an error if the execution fails or if the device is not found.
func (p *DeviceSetups) Delete(id models.DeviceID) error {
	slog.Debug("Delete device setup from database", "table", "device_setups", "id", id)

	query := `DELETE FROM device_setups WHERE device_id = ?`
	_, err := p.registry.db.Exec(query, id)
	return HandleSqlError(err)
}

// controlPins updates the pins on the device using the control package.
func (p *DeviceSetups) controlPins(deviceID models.DeviceID, pins []uint8) error {
	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return err
	}

	return control.SetPins(device, pins...)
}

// ScanDevicesSetup scans a database row into a DeviceSetup instance.
func ScanDevicesSetup(r Scannable) (*models.DeviceSetup, error) {
	deviceSetup := &models.DeviceSetup{}

	if err := r.Scan(&deviceSetup.DeviceID, &deviceSetup.Pins, &deviceSetup.CreatedAt); err != nil {
		return nil, err
	}

	return deviceSetup, nil
}
