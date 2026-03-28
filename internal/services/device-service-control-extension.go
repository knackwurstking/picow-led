package services

import (
	"fmt"
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
		if err == ErrorNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("get device for pins: %w", err)
	}

	pins, err := picow.GetPins(device)
	if err != nil {
		return nil, fmt.Errorf("get device pins from control layer: %w", err)
	}

	p.log.Debug("Pins for device %d: %v", deviceID, pins)

	// Cache the result
	p.pinCache.Store(deviceID, pins)

	return pins, nil
}

func (p *DeviceService) GetCurrentDuty(deviceID models.ID) ([]uint8, error) {
	// Get the device control record
	device, err := p.Get(deviceID)
	if err != nil {
		if err == ErrorNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("get device control for current duty: %w", err)
	}

	// Get the current duty from the picow device
	duty, err := picow.GetDuty(device)
	if err != nil {
		return nil, fmt.Errorf("get current duty from control layer: %w", err)
	}

	p.log.Debug("Current duty for device %d: %v", deviceID, duty)

	return duty, nil
}

func (p *DeviceService) GetVersion(deviceID models.ID) (string, error) {
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		if err == ErrorNotFound {
			return "", err
		}
		return "", fmt.Errorf("get device for version: %w", err)
	}

	version, err := picow.GetVersion(device)
	if err != nil {
		return "", fmt.Errorf("get version from control layer: %w", err)
	}

	return version, nil
}

func (p *DeviceService) GetDiskUsage(deviceID models.ID) (*picow.DiskUsage, error) {
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		if err == ErrorNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("get device for disk usage: %w", err)
	}

	diskUsage, err := picow.GetDiskUsage(device)
	if err != nil {
		return nil, fmt.Errorf("get disk usage from control layer: %w", err)
	}

	return diskUsage, nil
}

func (p *DeviceService) GetTemperature(deviceID models.ID) (float64, error) {
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		if err == ErrorNotFound {
			return 0, err
		}
		return 0, fmt.Errorf("get device for temperature: %w", err)
	}

	temperature, err := picow.GetTemperature(device)
	if err != nil {
		return 0, fmt.Errorf("get temperature from control layer: %w", err)
	}

	return temperature, nil
}

func (p *DeviceService) TogglePower(deviceID models.ID) ([]uint8, error) {
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		if err == ErrorNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("get device for power toggle: %w", err)
	}

	var newDuty []uint8
	currentDuty, err := picow.GetDuty(device)
	if err != nil {
		return nil, fmt.Errorf("get current duty for power toggle: %w", err)
	}

	if slices.Max(currentDuty) > 0 { // Just get the duty for turning OFF
		newDuty = make([]uint8, len(currentDuty))
		for i := range currentDuty {
			newDuty[i] = 0
		}
	} else { // Need to get the duty for turning ON (dc = deviceControl)
		if dc, _ := p.Get(device.ID); dc != nil && len(dc.Duty) > 0 {
			newDuty = dc.Duty // Get the duty from the database
		} else {
			// Nope, no duty in the database, get pins and set duty to 255 for each pin
			pins, err := p.GetPins(device.ID) // Use the cached version
			if err != nil {
				return nil, fmt.Errorf("get pins for power toggle: %w", err)
			}

			newDuty = make([]uint8, len(pins))
			for i := range pins {
				newDuty[i] = 255
			}
		}
	}

	p.log.Debug("Toggling power for device %d: current duty=%v, new duty=%v", deviceID, currentDuty, newDuty)

	if err = picow.SetDuty(device, newDuty...); err != nil {
		return nil, fmt.Errorf("set new duty for power toggle: %w", err)
	}

	return newDuty, nil
}

func (p *DeviceService) SetPins(deviceID models.ID, pins []uint8) error {
	if len(pins) == 0 {
		return ErrorValidation
	}

	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		if err == ErrorNotFound {
			return err
		}
		return fmt.Errorf("get device for set pins: %w", err)
	}

	p.log.Debug("Setting pins for device %d: %v", deviceID, pins)

	if err = picow.SetPins(device, pins...); err != nil {
		return fmt.Errorf("set pins in control layer: %w", err)
	}

	// Cache the new pins
	p.pinCache.Store(deviceID, pins)

	return nil
}

func (p *DeviceService) SetCurrentDuty(deviceID models.ID, duty []uint8) error {
	if len(duty) == 0 {
		return fmt.Errorf("duty cannot be empty")
	}

	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		if err == ErrorNotFound {
			return err
		}
		return fmt.Errorf("get device for set current duty: %w", err)
	}

	p.log.Debug("Setting current duty for device %d: %v", deviceID, duty)

	if err = picow.SetDuty(device, duty...); err != nil {
		return fmt.Errorf("set duty in control layer: %w", err)
	}

	// Store duty if bigger zero, else ignore
	if slices.Max(duty) > 0 {
		device.Duty = duty
		if err = p.Update(device); err != nil {
			return err
		}
	}

	return nil
}

func (p *DeviceService) TurnOn(deviceID models.ID) error {
	device, err := p.Get(deviceID)
	if err != nil {
		if err == ErrorNotFound {
			return err
		}
		return fmt.Errorf("get device for turn on: %w", err)
	}

	p.log.Debug("Turning on device %d with duty %v", deviceID, device.Duty)

	if err := picow.SetDuty(device, device.Duty...); err != nil {
		return fmt.Errorf("set duty in control layer for turn on: %w", err)
	}

	return nil
}

func (p *DeviceService) TurnOff(deviceID models.ID) error {
	device, err := p.registry.Device.Get(deviceID)
	if err != nil {
		if err == ErrorNotFound {
			return err
		}
		return fmt.Errorf("get device for turn off: %w", err)
	}

	// Get the current duty from the device
	currentDuty, err := picow.GetDuty(device)
	if err != nil {
		return fmt.Errorf("get current duty for turn off: %w", err)
	}

	// Set the duty to zero (turn off)
	duty := make([]uint8, len(currentDuty))

	p.log.Debug("Turning off device %d: current duty=%v, new duty=%v", deviceID, currentDuty, duty)

	if err := picow.SetDuty(device, duty...); err != nil {
		return fmt.Errorf("set duty in control layer for turn off: %w", err)
	}

	return nil
}
