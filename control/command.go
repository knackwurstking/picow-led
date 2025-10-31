package control

import (
	"encoding/json"
	"fmt"

	"github.com/knackwurstking/picow-led/env"
	"github.com/knackwurstking/picow-led/models"
)

// RunCommand sends a JSON-encoded request to the Picow device and reads the response.
func RunCommand[T any](id RequestID, device *models.ResolvedDevice, req any, resp *Response[T]) (T, error) {
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

// GetPins retrieves the current state of GPIO pins from the Picow device.
func GetPins(id RequestID, device *models.ResolvedDevice) ([]uint8, error) {
	req := NewGetPinsRequest(id)
	resp := &Response[[]uint8]{}
	return RunCommand(id, device, req, resp)
}

// SetPins sets the state of GPIO pins on the Picow device.
func SetPins(id RequestID, device *models.ResolvedDevice, pins ...uint8) error {
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

// GetColor retrieves the current color setting from the Picow device.
func GetColor(id RequestID, device *models.ResolvedDevice) ([]uint8, error) {
	req := NewGetColorRequest(id)
	resp := &Response[[]uint8]{}
	return RunCommand(id, device, req, resp)
}

// SetColor sets the color on the Picow device.
func SetColor(id RequestID, device *models.ResolvedDevice, duty ...uint8) error {
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

// GetTemperature retrieves the current temperature from the Picow device.
func GetTemperature(id RequestID, device *models.ResolvedDevice) (float32, error) {
	req := NewGetTemperatureRequest(id)
	resp := &Response[float32]{}
	return RunCommand(id, device, req, resp)
}

// GetDiskUsage retrieves disk usage information from the Picow device.
func GetDiskUsage(id RequestID, device *models.ResolvedDevice) (*DiskUsage, error) {
	req := NewGetDiskUsageRequest(id)
	resp := &Response[*DiskUsage]{}
	return RunCommand(id, device, req, resp)
}

// GetVersion retrieves the current version of the firmware running on the Picow device.
func GetVersion(id RequestID, device *models.ResolvedDevice) (string, error) {
	req := NewGetVersionRequest(id)
	resp := &Response[string]{}
	return RunCommand(id, device, req, resp)
}
