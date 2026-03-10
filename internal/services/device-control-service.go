package services

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/knackwurstking/picow-led/internal/models"
)

type DeviceControlService struct {
	registry *Registry

	// Cache for device pins
	pinCache sync.Map
}

func NewDeviceControlService(registry *Registry) *DeviceControlService {
	return &DeviceControlService{
		registry: registry,
	}
}

// TODO: Remove this table... Make this type an extension of the *Device type for control specific methods
func (p *DeviceControlService) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS device_controls (
		device_id INTEGER NOT NULL,
		color TEXT NOT NULL,

		PRIMARY KEY ("device_id")
	);`
	_, err := p.registry.db.Exec(query)
	if err != nil {
		return NewServiceError("create device controls table", err)
	}
	return nil
}

func (p *DeviceControlService) Get(deviceID models.ID) (*models.DeviceControl, error) {
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

func (p *DeviceControlService) List() ([]*models.DeviceControl, error) {
	query := `SELECT * FROM device_controls`
	rows, err := p.registry.db.Query(query)
	if err != nil {
		return nil, NewServiceError("list device controls", HandleSqlError(err))
	}
	defer rows.Close()

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

func (p *DeviceControlService) Add(deviceControl *models.DeviceControl) (models.ID, error) {
	if deviceControl.Validate() != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidDeviceSetup, "device control validation failed")
	}

	query := `INSERT INTO device_controls (device_id, color) VALUES (?, ?)`
	result, err := p.registry.db.Exec(query, deviceControl.DeviceID, deviceControl.Color)
	if err != nil {
		return 0, NewServiceError("add device control", HandleSqlError(err))
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, NewServiceError("get last inserted device control ID", HandleSqlError(err))
	}

	return models.ID(id), nil
}

func (p *DeviceControlService) Update(deviceControl *models.DeviceControl) error {
	if deviceControl.Validate() != nil {
		return fmt.Errorf("%w: %v", ErrInvalidDeviceSetup, "device control validation failed")
	}

	query := `UPDATE device_controls SET color = ? WHERE device_id = ?`
	_, err := p.registry.db.Exec(query, deviceControl.Color, deviceControl.DeviceID)
	if err != nil {
		return NewServiceError("update device control", HandleSqlError(err))
	}

	return nil
}

func (p *DeviceControlService) Delete(deviceID models.ID) error {
	query := `DELETE FROM device_controls WHERE device_id = ?`
	_, err := p.registry.db.Exec(query, deviceID)
	if err != nil {
		return NewServiceError("delete device control", HandleSqlError(err))
	}

	return nil
}

func (p *DeviceControlService) GetPins(deviceID models.ID) ([]uint8, error) {
	// Check if we have a cached value
	if cached, ok := p.pinCache.Load(deviceID); ok {
		return cached.([]uint8), nil
	}

	// If not cached, fetch from control layer
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		return nil, NewServiceError("get device for pins", err)
	}

	pins, err := control.GetPins(device)
	if err != nil {
		return nil, NewServiceError("get device pins from control layer", err)
	}

	// Cache the result
	p.pinCache.Store(deviceID, pins)

	return pins, nil
}

func (p *DeviceControlService) GetCurrentColor(deviceID models.ID) ([]uint8, error) {
	// Fetch device
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		return nil, NewServiceError("get device for current color", err)
	}

	// Get the device control record
	deviceControl, err := p.Get(device.ID)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result") {
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

func (p *DeviceControlService) GetVersion(deviceID models.ID) (string, error) {
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		return "", NewServiceError("get device for version", err)
	}

	version, err := control.GetVersion(device)
	if err != nil {
		return "", NewServiceError("get version from control layer", err)
	}

	return version, nil
}

func (p *DeviceControlService) GetDiskUsage(deviceID models.ID) (*control.DiskUsage, error) {
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		return nil, NewServiceError("get device for disk usage", err)
	}

	diskUsage, err := control.GetDiskUsage(device)
	if err != nil {
		return nil, NewServiceError("get disk usage from control layer", err)
	}

	return diskUsage, nil
}

func (p *DeviceControlService) GetTemperature(deviceID models.ID) (float64, error) {
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		return 0, NewServiceError("get device for temperature", err)
	}

	temperature, err := control.GetTemperature(device)
	if err != nil {
		return 0, NewServiceError("get temperature from control layer", err)
	}

	return temperature, nil
}

func (p *DeviceControlService) TogglePower(deviceID models.ID) ([]uint8, error) {
	device, err := p.registry.Device.Get(deviceID)
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

	if err = control.SetColor(device, newColor...); err != nil {
		return nil, NewServiceError("set color for power toggle", err)
	}

	return newColor, nil
}

func (p *DeviceControlService) SetCurrentColor(deviceID models.ID, color []uint8) error {
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		return NewServiceError("get device for setting current color", err)
	}

	if slices.Max(color) > 0 {
		deviceControl := models.NewDeviceControl(deviceID, color...)
		if err = p.Update(deviceControl); err != nil {
			if IsNotFoundError(err) {
				if err = p.setInitialEntry(deviceID); err != nil {
					return NewServiceError("set initial entry for setting current color", err)
				}

				deviceControl := models.NewDeviceControl(deviceID, color...)
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

func (p *DeviceControlService) TurnOn(deviceID models.ID) error {
	device, err := p.registry.Device.Get(deviceID)
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

func (p *DeviceControlService) TurnOff(deviceID models.ID) error {
	device, err := p.registry.Device.Get(deviceID)
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

func (p *DeviceControlService) setInitialEntry(deviceID models.ID) error {
	// Get pins using the cached version
	pins, err := p.GetPins(deviceID)
	if err != nil {
		return NewServiceError("get pins for initial entry", err)
	}

	color := make([]uint8, len(pins))
	for i := range pins {
		color[i] = 255
	}

	data := models.NewDeviceControl(deviceID, color...)
	if _, err := p.Add(data); err != nil {
		return NewServiceError("add initial device control entry", err)
	}

	return nil
}

func ScanDeviceControl(scanner Scannable) (*models.Device, error) {
	control := &models.Device{}
	err := scanner.Scan(&control.ID, &control.Color)
	if err != nil {
		return nil, NewServiceError("scan device control", err)
	}
	return control, nil
}

var _ Service = (*DeviceControlService)(nil)
