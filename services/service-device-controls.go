package services

import (
	"database/sql"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/knackwurstking/picow-led/control"
	"github.com/knackwurstking/picow-led/models"
)

// DeviceControls struct manages the control of devices.
type DeviceControls struct {
	registry *Registry

	// Cache for device pins
	pinCache sync.Map
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
	if err != nil {
		return NewServiceError("create device controls table", err)
	}
	return nil
}

// Get retrieves a device control record from the database by device ID.
// It returns a DeviceControl instance and an error if the record is not found or execution fails.
func (p *DeviceControls) Get(deviceID models.DeviceID) (*models.DeviceControl, error) {
	slog.Debug("Get device control record", "id", deviceID)

	query := `SELECT * FROM device_controls WHERE device_id = ?`
	row := p.registry.db.QueryRow(query, deviceID)
	deviceControl, err := ScanDeviceControl(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError(fmt.Sprintf("deviceID %d", deviceID))
		}
		return nil, NewServiceError("get device control", HandleSqlError(err))
	}
	return deviceControl, nil
}

// List retrieves all device control records from the database.
// It returns a slice of DeviceControl instances and an error if execution fails.
func (p *DeviceControls) List() ([]*models.DeviceControl, error) {
	slog.Debug("Get all device control records")

	query := `SELECT * FROM device_controls`
	rows, err := p.registry.db.Query(query)
	if err != nil {
		return nil, NewServiceError("list device controls", HandleSqlError(err))
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Warn("close device controls rows", "error", err)
		}
	}()

	deviceControls := make([]*models.DeviceControl, 0)
	for rows.Next() {
		deviceControl, err := ScanDeviceControl(rows)
		if err != nil {
			return nil, NewServiceError("scan device control from rows", HandleSqlError(err))
		}
		deviceControls = append(deviceControls, deviceControl)
	}

	if err := rows.Err(); err != nil {
		return nil, NewServiceError("iterate device control rows", err)
	}

	return deviceControls, nil
}

// Add inserts a new device control record into the database.
// It returns the inserted device ID and an error if execution fails.
func (p *DeviceControls) Add(deviceControl *models.DeviceControl) (models.DeviceID, error) {
	if !deviceControl.Validate() {
		return 0, fmt.Errorf("%w: %v", ErrInvalidDeviceSetup, "device control validation failed")
	}

	slog.Debug("Adding device control", "id", deviceControl.DeviceID,
		"color", deviceControl.Color, "modified", deviceControl.ModifiedAt)

	query := `INSERT INTO device_controls (device_id, color) VALUES (?, ?)`
	result, err := p.registry.db.Exec(query, deviceControl.DeviceID, deviceControl.Color)
	if err != nil {
		return 0, NewServiceError("add device control", HandleSqlError(err))
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, NewServiceError("get last inserted device control ID", HandleSqlError(err))
	}

	return models.DeviceID(id), nil
}

// Update updates an existing device control record in the database.
// It returns an error if execution fails or the record is not found.
func (p *DeviceControls) Update(deviceControl *models.DeviceControl) error {
	if !deviceControl.Validate() {
		return fmt.Errorf("%w: %v", ErrInvalidDeviceSetup, "device control validation failed")
	}
	deviceControl.ModifiedAt = time.Now()

	slog.Debug("Updating device control", "id", deviceControl.DeviceID,
		"color", deviceControl.Color, "modified", deviceControl.ModifiedAt)

	query := `UPDATE device_controls SET color = ? WHERE device_id = ?`
	_, err := p.registry.db.Exec(query, deviceControl.Color, deviceControl.DeviceID)
	if err != nil {
		return NewServiceError("update device control", HandleSqlError(err))
	}

	return nil
}

// Delete removes the record associated with a given device ID from the database.
// It returns an error if the execution fails or if the device is not found.
func (p *DeviceControls) Delete(deviceID models.DeviceID) error {
	slog.Debug("Deleting device control", "id", deviceID)

	query := `DELETE FROM device_controls WHERE device_id = ?`
	_, err := p.registry.db.Exec(query, deviceID)
	if err != nil {
		return NewServiceError("delete device control", HandleSqlError(err))
	}

	return nil
}

// GetPins retrieves the pins for a device, with caching.
func (p *DeviceControls) GetPins(deviceID models.DeviceID) ([]uint8, error) {
	slog.Debug("Getting device pins", "id", deviceID)

	// Check if we have a cached value
	if cached, ok := p.pinCache.Load(deviceID); ok {
		slog.Debug("Using cached pins", "id", deviceID)
		return cached.([]uint8), nil
	}

	// If not cached, fetch from control layer
	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return nil, NewServiceError("get device for pins", err)
	}

	slog.Debug("Running GetPins from the controls package", "id", deviceID)
	pins, err := control.GetPins(device)
	if err != nil {
		return nil, NewServiceError("get device pins from control layer", err)
	}

	// Cache the result
	p.pinCache.Store(deviceID, pins)

	slog.Debug("Cached pins", "id", deviceID)
	return pins, nil
}

// GetCurrentColor retrieves the current color of the device from the database and will auto update the database if the color is different from the stored color and not 0.
func (p *DeviceControls) GetCurrentColor(deviceID models.DeviceID) ([]uint8, error) {
	slog.Debug("Getting device current color", "id", deviceID)

	// Fetch device
	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return nil, NewServiceError("get device for current color", err)
	}

	// Get the device control record
	deviceControl, err := p.Get(device.ID)
	if err != nil {
		if !IsNotFoundError(err) {
			return nil, NewServiceError("get device control for current color", err)
		}

		if err = p.setInitialEntry(device.ID); err != nil {
			return nil, NewServiceError("set initial entry for current color", err)
		}

		deviceControl, err = p.Get(device.ID)
		if err != nil {
			return nil, NewServiceError("get device control after setting initial entry", err)
		}
	}

	// Get the current color from the picow device
	color, err := control.GetColor(device)
	if err != nil {
		return nil, NewServiceError("get color from device", err)
	}

	// Check if the color is different from the stored color and not 0
	if slices.Max(color) > 0 {
		// Update the device control object with the new color
		deviceControl.Color = color
		if err = p.Update(deviceControl); err != nil {
			return nil, NewServiceError("update device control with new color", err)
		}
	}

	return color, nil
}

func (p *DeviceControls) GetVersion(deviceID models.DeviceID) (string, error) {
	slog.Debug("Get version", "id", deviceID)

	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return "", NewServiceError("get device for version", err)
	}

	version, err := control.GetVersion(device)
	if err != nil {
		return "", NewServiceError("get version from control layer", err)
	}

	return version, nil
}

func (p *DeviceControls) GetDiskUsage(deviceID models.DeviceID) (*control.DiskUsage, error) {
	slog.Debug("Get disk usage", "id", deviceID)

	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return nil, NewServiceError("get device for disk usage", err)
	}

	diskUsage, err := control.GetDiskUsage(device)
	if err != nil {
		return nil, NewServiceError("get disk usage from control layer", err)
	}

	return diskUsage, nil
}

func (p *DeviceControls) GetTemperature(deviceID models.DeviceID) (float64, error) {
	slog.Debug("Get temperature", "id", deviceID)

	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return 0, NewServiceError("get device for temperature", err)
	}

	temperature, err := control.GetTemperature(device)
	if err != nil {
		return 0, NewServiceError("get temperature from control layer", err)
	}

	return temperature, nil
}

func (p *DeviceControls) TogglePower(deviceID models.DeviceID) ([]uint8, error) {
	slog.Debug("Toggle power", "id", deviceID)

	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return nil, NewServiceError("get device for power toggle", err)
	}

	var newColor []uint8
	currentColor, err := control.GetColor(device)
	if err != nil {
		return nil, NewServiceError("get current color for power toggle", err)
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
			pins, err := p.GetPins(device.ID) // Use the cached version
			if err != nil {
				return nil, NewServiceError("get pins for power toggle", err)
			}

			newColor = make([]uint8, len(pins))
			for i := range pins {
				newColor[i] = 255
			}
		}
	}

	err = control.SetColor(device, newColor...)
	if err != nil {
		return nil, NewServiceError("set color for power toggle", err)
	}

	return newColor, nil
}

func (p *DeviceControls) SetCurrentColor(deviceID models.DeviceID, color []uint8) error {
	slog.Debug("Set current color", "id", deviceID)

	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return NewServiceError("get device for setting current color", err)
	}

	if slices.Max(color) > 0 {
		deviceControl := models.NewDeviceControl(deviceID, color)
		if err = p.Update(deviceControl); err != nil {
			if IsNotFoundError(err) {
				if err = p.setInitialEntry(deviceID); err != nil {
					return NewServiceError("set initial entry for setting current color", err)
				}

				deviceControl := models.NewDeviceControl(deviceID, color)
				if err = p.Update(deviceControl); err != nil {
					return NewServiceError("update device control after setting initial entry", err)
				}
			} else {
				return NewServiceError("update device control for setting current color", err)
			}
		}
	}

	if err := control.SetColor(device, color...); err != nil {
		return NewServiceError("set color in control layer", err)
	}

	return nil
}

func (p *DeviceControls) TurnOn(deviceID models.DeviceID) error {
	slog.Debug("Turn on", "id", deviceID)

	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return NewServiceError("get device for turn on", err)
	}

	deviceControl, err := p.Get(deviceID)
	if err != nil {
		if IsNotFoundError(err) {
			if err = p.setInitialEntry(deviceID); err != nil {
				return NewServiceError("set initial entry for turn on", err)
			}

			deviceControl, err = p.Get(deviceID)
			if err != nil {
				return NewServiceError("get device control after setting initial entry for turn on", err)
			}
		} else {
			return NewServiceError("get device control for turn on", err)
		}
	}

	if err := control.SetColor(device, deviceControl.Color...); err != nil {
		return NewServiceError("set color in control layer for turn on", err)
	}

	return nil
}

// TurnOff turns off the device by setting its color to zero.
func (p *DeviceControls) TurnOff(deviceID models.DeviceID) error {
	slog.Debug("Turn off", "id", deviceID)

	device, err := p.registry.Devices.Get(deviceID)
	if err != nil {
		return NewServiceError("get device for turn off", err)
	}

	// Get the current color from the device
	currentColor, err := control.GetColor(device)
	if err != nil {
		return NewServiceError("get current color for turn off", err)
	}

	// Set the color to zero (turn off)
	color := make([]uint8, len(currentColor))
	if err := control.SetColor(device, color...); err != nil {
		return NewServiceError("set color in control layer for turn off", err)
	}

	return nil
}

func (p *DeviceControls) setInitialEntry(deviceID models.DeviceID) error {
	slog.Debug("Set the initial entry", "id", deviceID)

	// Get pins using the cached version
	pins, err := p.GetPins(deviceID)
	if err != nil {
		return NewServiceError("get pins for initial entry", err)
	}

	color := make([]uint8, len(pins))
	for i := range pins {
		color[i] = 255
	}

	data := models.NewDeviceControl(deviceID, color)
	if _, err := p.Add(data); err != nil {
		return NewServiceError("add initial device control entry", err)
	}

	return nil
}

func ScanDeviceControl(scanner Scannable) (*models.DeviceControl, error) {
	control := &models.DeviceControl{}
	err := scanner.Scan(&control.DeviceID, &control.Color, &control.ModifiedAt, &control.CreatedAt)
	if err != nil {
		return nil, NewServiceError("scan device control", err)
	}
	return control, nil
}
