package services

import (
	"slices"

	"github.com/knackwurstking/picow-led/pkg/models"
	"github.com/knackwurstking/picow-led/pkg/picow"
)

func (p *DeviceService) GetPins(deviceID models.ID) ([]uint8, error) {
	// Check if we have a cached value
	if cached, ok := p.pinCache.Load(deviceID); ok {
		return cached.([]uint8), nil
	}

	// If not cached, fetch from control layer
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		return nil, NewServiceError("get device for pins", err)
	}

	pins, err := picow.GetPins(device)
	if err != nil {
		return nil, NewServiceError("get device pins from control layer", err)
	}

	// Cache the result
	p.pinCache.Store(deviceID, pins)

	return pins, nil
}

func (p *DeviceService) GetCurrentColor(deviceID models.ID) ([]uint8, error) {
	// Get the device control record
	device, err := p.Get(deviceID)
	if err != nil {
		return nil, NewServiceError("get device control for current color", err)
	}

	// Get the current color from the picow device
	color, err := picow.GetColor(device)
	if err != nil {
		return nil, NewServiceError("get color from device", err)
	}

	return color, nil
}

func (p *DeviceService) GetVersion(deviceID models.ID) (string, error) {
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		return "", NewServiceError("get device for version", err)
	}

	version, err := picow.GetVersion(device)
	if err != nil {
		return "", NewServiceError("get version from control layer", err)
	}

	return version, nil
}

func (p *DeviceService) GetDiskUsage(deviceID models.ID) (*picow.DiskUsage, error) {
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		return nil, NewServiceError("get device for disk usage", err)
	}

	diskUsage, err := picow.GetDiskUsage(device)
	if err != nil {
		return nil, NewServiceError("get disk usage from control layer", err)
	}

	return diskUsage, nil
}

func (p *DeviceService) GetTemperature(deviceID models.ID) (float64, error) {
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		return 0, NewServiceError("get device for temperature", err)
	}

	temperature, err := picow.GetTemperature(device)
	if err != nil {
		return 0, NewServiceError("get temperature from control layer", err)
	}

	return temperature, nil
}

func (p *DeviceService) TogglePower(deviceID models.ID) ([]uint8, error) {
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		return nil, NewServiceError("get device for power toggle", err)
	}

	var newColor []uint8
	currentColor, err := picow.GetColor(device)
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

	if err = picow.SetColor(device, newColor...); err != nil {
		return nil, NewServiceError("set color for power toggle", err)
	}

	return newColor, nil
}

func (p *DeviceService) SetCurrentColor(deviceID models.ID, color []uint8) error {
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		return NewServiceError("get device for setting current color", err)
	}

	if err := picow.SetColor(device, color...); err != nil {
		return NewServiceError("set color in control layer", err)
	}

	return nil
}

func (p *DeviceService) TurnOn(deviceID models.ID) error {
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		return NewServiceError("get device for turn on", err)
	}

	deviceControl, err := p.Get(deviceID)
	if err != nil {
		return NewServiceError("get device control for turn on", err)
	}

	if err := picow.SetColor(device, deviceControl.Color...); err != nil {
		return NewServiceError("set color in control layer for turn on", err)
	}

	return nil
}

func (p *DeviceService) TurnOff(deviceID models.ID) error {
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		return NewServiceError("get device for turn off", err)
	}

	// Get the current color from the device
	currentColor, err := picow.GetColor(device)
	if err != nil {
		return NewServiceError("get current color for turn off", err)
	}

	// Set the color to zero (turn off)
	color := make([]uint8, len(currentColor))
	if err := picow.SetColor(device, color...); err != nil {
		return NewServiceError("set color in control layer for turn off", err)
	}

	return nil
}
