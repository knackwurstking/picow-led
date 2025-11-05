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
	slog.Debug("Get device control record", "id", deviceID)

	query := `SELECT * FROM device_controls WHERE device_id = ?`
	row := p.registry.db.QueryRow(query, deviceID)
	deviceControl, err := ScanDeviceControl(row)
	return deviceControl, HandleSqlError(err)
}

// List retrieves all device control records from the database.
// It returns a slice of DeviceControl instances and an error if execution fails.
func (p *DeviceControls) List() ([]*models.DeviceControl, error) {
	slog.Debug("Get all device control records")

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
	if !deviceControl.Validate() {
		return 0, ErrInvalidDeviceSetup
	}

	slog.Debug("Adding device control", "id", deviceControl.DeviceID,
		"color", deviceControl.Color, "modified", deviceControl.ModifiedAt)

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
	if !deviceControl.Validate() {
		return ErrInvalidDeviceSetup
	}

	slog.Debug("Updating device control", "id", deviceControl.DeviceID,
		"color", deviceControl.Color, "modified", deviceControl.ModifiedAt)

	query := `UPDATE device_controls SET color = ? WHERE device_id = ?`
	_, err := p.registry.db.Exec(query, deviceControl.Color, deviceControl.DeviceID)
	return HandleSqlError(err)
}

// Delete removes the record associated with a given device ID from the database.
// It returns an error if the execution fails or if the device is not found.
func (p *DeviceControls) Delete(deviceID models.DeviceID) error {
	slog.Debug("Deleting device control", "id", deviceID)

	query := `DELETE FROM device_controls WHERE device_id = ?`
	_, err := p.registry.db.Exec(query, deviceID)
	return HandleSqlError(err)
}

func (p *DeviceControls) GetPins(deviceID models.DeviceID) ([]uint8, error) {
	slog.Debug("Getting device pins", "id", deviceID)

	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		if IsNotFoundError(err) {
			return nil, fmt.Errorf("device %d not found", deviceID)
		}
		return nil, err
	}

	slog.Debug("Running GetPins from the controls package", "id", deviceID)
	pins, err := control.GetPins(device)
	if err != nil {
		return nil, fmt.Errorf("failed to get device pins: %v", err)
	}

	return pins, nil
}

// CurrentColor retrieves the current color of the device from the database and will auto update the database if the color is different from the stored color and not 0.
func (p *DeviceControls) GetCurrentColor(deviceID models.DeviceID) ([]uint8, error) {
	slog.Debug("Getting device current color", "id", deviceID)

	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %v", err)
	}

	deviceControl, err := p.Get(device.ID)
	if err != nil {
		if !IsNotFoundError(err) {
			return nil, fmt.Errorf("failed to get device control: %v", err)
		}

		if err = p.setInitialEntry(device.ID); err != nil {
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
	slog.Debug("Get version", "id", deviceID)

	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return "", err
	}

	return control.GetVersion(device)
}

func (p *DeviceControls) GetDiskUsage(deviceID models.DeviceID) (*control.DiskUsage, error) {
	slog.Debug("Get disk usage", "id", deviceID)

	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return nil, err
	}

	return control.GetDiskUsage(device)
}

func (p *DeviceControls) GetTemperature(deviceID models.DeviceID) (float64, error) {
	slog.Debug("Get temperature", "id", deviceID)

	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return 0, err
	}

	return control.GetTemperature(device)
}

func (p *DeviceControls) TogglePower(deviceID models.DeviceID) ([]uint8, error) {
	slog.Debug("Toggle power", "id", deviceID)

	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return nil, err
	}

	var newColor []uint8
	currentColor, err := control.GetColor(device)
	if err != nil {
		return nil, err
	}

	if slices.Max(currentColor) > 0 { // Just get the color for turning OFF
		newColor = make([]uint8, len(currentColor))
		for i := range currentColor {
			newColor[i] = 0
		}
	} else { // Need to get the color for turning ON (dc = deviceControl)
		if dc, _ := p.Get(device.ID); dc != nil && len(dc.Color) > 0 {
			newColor = dc.Color // Get the color from the database
		} else {
			// Nope, no color in the database, get pins and set color to 255 for each pin
			pins, err := control.GetPins(device)
			if err != nil {
				return nil, err
			}

			newColor = make([]uint8, len(pins))
			for i := range pins {
				newColor[i] = 255
			}
		}
	}

	return newColor, control.SetColor(device, newColor...)
}

func (p *DeviceControls) setInitialEntry(deviceID models.DeviceID) error {
	slog.Debug("Set the initial entry", "id", deviceID)

	pins, err := p.GetPins(deviceID)
	if err != nil {
		return err
	}

	color := make([]uint8, len(pins))
	for i, _ := range pins {
		color[i] = 255
	}

	data := models.NewDeviceControl(deviceID, color)
	if _, err := p.Add(data); err != nil {
		return err
	}

	return nil
}

func ScanDeviceControl(scanner Scannable) (*models.DeviceControl, error) {
	control := &models.DeviceControl{}
	err := scanner.Scan(&control.DeviceID, &control.Color, &control.ModifiedAt, &control.CreatedAt)
	return control, err
}
