package services

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/knackwurstking/picow-led/control"
	"github.com/knackwurstking/picow-led/models"
)

// DeviceControls struct manages the control of devices.
type DeviceControls struct {
	registry *Registry
}

// NewDeviceControls creates a new instance of DeviceControls with the provided registry.
func NewDeviceControls(registry *Registry) *DeviceControls {
	return &DeviceControls{
		registry: registry,
	}
}

// CreateTable creates a table for device control data if it doesn't already exist.
// It returns an error if the database execution fails.
func (p *DeviceControls) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS device_controls (
		device_id INTEGER PRIMARY KEY NOT NULL,
		color TEXT NOT NULL,
		modified_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := p.registry.db.Exec(query)
	return err
}

// Get retrieves a device control record from the database by device ID.
// It returns a DeviceControl instance and an error if the record is not found or execution fails.
func (p *DeviceControls) Get(deviceID models.DeviceID) (*models.DeviceControl, error) {
	slog.Debug("Get device control from database", "table", "device_controls", "id", deviceID)

	query := `SELECT * FROM device_controls WHERE device_id = ?`
	row := p.registry.db.QueryRow(query, deviceID)
	deviceControl, err := ScanDeviceControl(row)
	return deviceControl, HandleSqlError(err)
}

// List retrieves all device control records from the database.
// It returns a slice of DeviceControl instances and an error if execution fails.
func (p *DeviceControls) List() ([]*models.DeviceControl, error) {
	slog.Debug("List device controls from database", "table", "device_controls")

	query := `SELECT * FROM device_controls`
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
func (p *DeviceControls) Add(deviceControl *models.DeviceControl) (models.DeviceID, error) {
	slog.Debug("Add device control to database", "table", "device_controls", "device_id", deviceControl.DeviceID)

	query := `INSERT INTO device_controls (device_id, color) VALUES (?, ?)`
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
func (p *DeviceControls) Update(deviceControl *models.DeviceControl) error {
	slog.Debug("Update device control in database", "table", "device_controls", "id", deviceControl.DeviceID)

	query := `UPDATE device_controls SET color = ? WHERE device_id = ?`
	_, err := p.registry.db.Exec(query, deviceControl.Color, deviceControl.DeviceID)
	return HandleSqlError(err)
}

// Delete removes the record associated with a given device ID from the database.
// It returns an error if the execution fails or if the device is not found.
func (p *DeviceControls) Delete(deviceID models.DeviceID) error {
	slog.Debug("Delete device control from database", "table", "device_controls", "id", deviceID)

	query := `DELETE FROM device_controls WHERE device_id = ?`
	_, err := p.registry.db.Exec(query, deviceID)
	return HandleSqlError(err)
}

func (p *DeviceControls) GetPins(deviceID models.DeviceID) ([]uint8, error) {
	slog.Debug("Get device pins from database", "table", "device_controls", "id", deviceID)

	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return nil, err
	}

	pins, err := control.GetPins(device)
	if err != nil {
		return nil, err
	}

	// Store the pins in the database
	_, err = p.registry.db.Exec("UPDATE device_controls SET pins = ? WHERE device_id = ?", pins, deviceID)
	if err != nil {
		return nil, HandleSqlError(err)
	}

	return pins, nil
}

// CurrentColor retrieves the current color of the device from the database and will auto update the database if the color is different from the stored color and not 0.
func (p *DeviceControls) GetCurrentColor(deviceID models.DeviceID) ([]uint8, error) {
	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return nil, err
	}

	deviceControl, err := p.Get(device.ID)
	if err != nil {
		if err != ErrNotFound {
			return nil, err
		}

		pins, err := p.GetPins(device.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get pins for device %d: %v", device.ID, err)
		}

		// Create the initial color (duty) for each pin
		initialColor := make([]uint8, len(pins))
		for i := range initialColor {
			initialColor[i] = 255
		}

		deviceControl = models.NewDeviceControl(device.ID, initialColor)
		if _, err := p.Add(deviceControl); err != nil {
			return nil, err
		}
	}

	// Get the current color from the picow device
	color, err := control.GetColor(device)
	if err != nil {
		return nil, err
	}
	// Check if the color is different from the stored color and not 0
	if slices.Max(color) > 0 {
		// Update the device control object with the new color
		deviceControl.Color = color
		if err = p.Update(deviceControl); err != nil {
			return nil, err
		}
	}

	return color, nil
}

func (p *DeviceControls) GetVersion(deviceID models.DeviceID) (string, error) {
	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return "", err
	}

	return control.GetVersion(device)
}

func (p *DeviceControls) GetDiskUsage(deviceID models.DeviceID) (*control.DiskUsage, error) {
	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return nil, err
	}

	return control.GetDiskUsage(device)
}

func (p *DeviceControls) GetTemperature(deviceID models.DeviceID) (float64, error) {
	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return 0, err
	}

	return control.GetTemperature(device)
}

func ScanDeviceControl(scanner Scannable) (*models.DeviceControl, error) {
	control := &models.DeviceControl{}
	err := scanner.Scan(&control.DeviceID, &control.Color, &control.ModifiedAt, &control.CreatedAt)
	return control, err
}
