package control

import (
	"encoding/json"
	"fmt"

	"github.com/knackwurstking/picow-led/env"
	"github.com/knackwurstking/picow-led/models"
)

// RunCommand sends a JSON-encoded request to the Picow device and reads the response.
func RunCommand[T any](id RequestID, device *models.Device, req any, resp *Response[T]) (T, error) {
	picow := NewPicoW(device)

	data, err := json.Marshal(req)
	if err != nil {
		return *new(T), fmt.Errorf("failed to marshal request: %v", err)
	}
	n, err := picow.Write(data)
	if err != nil {
		return *new(T), err
	}
	if n != len(data) {
		return *new(T), ErrNoData
	}

	if id == RequestIDNoResponse {
		return *new(T), nil
	}

	respData, err := picow.ReadAll()
	if err != nil {
		return *new(T), fmt.Errorf("failed to read response: %v", err)
	}
	if err = json.Unmarshal(respData, &resp); err != nil {
		return *new(T), fmt.Errorf("failed to unmarshal response: %v", err)
	}
	if resp.Error != "" {
		return *new(T), fmt.Errorf("picow error: %s", resp.Error)
	}

	return resp.Data, nil
}

// GetPins retrieves the current state of GPIO pins from the Picow device. Default Request ID will be used.
func GetPins(device *models.Device) ([]uint8, error) {
	return GetPinsWithID(RequestIDDefault, device)
}

// GetPinsWithID retrieves the current state of GPIO pins from the Picow device.
func GetPinsWithID(id RequestID, device *models.Device) ([]uint8, error) {
	req := NewGetPinsRequest(id)
	resp := &Response[[]uint8]{}
	return RunCommand(id, device, req, resp)
}

// SetPins sets the state of GPIO pins on the Picow device. Default Request ID will be used.
func SetPins(device *models.Device, pins ...uint8) error {
	return SetPinsWithID(RequestIDDefault, device, pins...)
}

// SetPinsWithID sets the state of GPIO pins on the Picow device.
func SetPinsWithID(id RequestID, device *models.Device, pins ...uint8) error {
	// Validate pins first
	for _, pin := range pins {
		if pin < env.MinPin || pin > env.MaxPin {
			return fmt.Errorf("invalid pin: %d (min: %d, max: %d)", pin, env.MinPin, env.MaxPin)
		}
	}

	req := NewSetPinsRequest(id, pins)
	resp := &Response[struct{}]{}
	_, err := RunCommand(id, device, req, resp)
	return err
}

// GetColor retrieves the current color setting from the Picow device. Default Request ID will be used.
func GetColor(device *models.Device) ([]uint8, error) {
	return GetColorWithID(RequestIDDefault, device)
}

// GetColorWithID retrieves the current color setting from the Picow device.
func GetColorWithID(id RequestID, device *models.Device) ([]uint8, error) {
	req := NewGetColorRequest(id)
	resp := &Response[[]uint8]{}
	return RunCommand(id, device, req, resp)
}

// SetColor sets the color on the Picow device. Default request ID will be used.
func SetColor(device *models.Device, duty ...uint8) error {
	return SetColorWithID(RequestIDDefault, device, duty...)
}

// SetColorWithID sets the color on the Picow device.
func SetColorWithID(id RequestID, device *models.Device, duty ...uint8) error {
	// Validate duty first
	for _, d := range duty {
		if d < env.MinDuty || d > env.MaxDuty {
			return fmt.Errorf("invalid duty: %d (min: %d, max: %d)", d, env.MinDuty, env.MaxDuty)
		}
	}

	req := NewSetColorRequest(id, duty)
	resp := &Response[struct{}]{}
	_, err := RunCommand(id, device, req, resp)
	return err
}

// GetTemperature retrieves the current temperature from the Picow device. Default request ID will be used.
func GetTemperature(device *models.Device) (float32, error) {
	return GetTemperatureWithID(RequestIDDefault, device)
}

// GetTemperatureWithID retrieves the current temperature from the Picow device.
func GetTemperatureWithID(id RequestID, device *models.Device) (float32, error) {
	req := NewGetTemperatureRequest(id)
	resp := &Response[float32]{}
	return RunCommand(id, device, req, resp)
}

// GetDiskUsage retrieves disk usage information from the Picow device. Default request ID will be used.
func GetDiskUsage(device *models.Device) (*DiskUsage, error) {
	return GetDiskUsageWithID(RequestIDDefault, device)
}

// GetDiskUsageWithID retrieves disk usage information from the Picow device.
func GetDiskUsageWithID(id RequestID, device *models.Device) (*DiskUsage, error) {
	req := NewGetDiskUsageRequest(id)
	resp := &Response[*DiskUsage]{}
	return RunCommand(id, device, req, resp)
}

// GetVersion retrieves the current version of the firmware running on the Picow device. Default request ID will be used.
func GetVersion(device *models.Device) (string, error) {
	return GetVersionWithID(RequestIDDefault, device)
}

// GetVersionWithID retrieves the current version of the firmware running on the Picow device.
func GetVersionWithID(id RequestID, device *models.Device) (string, error) {
	req := NewGetVersionRequest(id)
	resp := &Response[string]{}
	return RunCommand(id, device, req, resp)
}
