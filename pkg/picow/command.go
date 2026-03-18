package picow

import (
	"encoding/json"
	"fmt"

	"github.com/knackwurstking/picow-led/pkg/models"
)

// RunCommand sends a JSON-encoded request to the Picow device and reads the response.
func RunCommand[T any](id RequestID, device *models.Device, req any, resp *Response[T]) (T, error) {
	picow := NewPicoW(device)

	if err := picow.Connect(); err != nil {
		return *new(T), fmt.Errorf("connect: %v", err)
	}
	defer picow.Close()

	data, err := json.Marshal(req)
	if err != nil {
		return *new(T), fmt.Errorf("marshal request: %v", err)
	}
	_, err = picow.Write(data)
	if err != nil {
		return *new(T), fmt.Errorf("write request: %v", err)
	}

	if id == RequestIDNoResponse {
		return *new(T), nil
	}

	respData, err := picow.ReadAll()
	if err != nil {
		return *new(T), fmt.Errorf("response: %v", err)
	}
	if err = json.Unmarshal(respData, &resp); err != nil {
		return *new(T), fmt.Errorf("unmarshal response: %v", err)
	}
	if resp.Error != "" {
		return *new(T), fmt.Errorf("device responded with error: %s", resp.Error)
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
		if pin < MinPin || pin > MaxPin {
			return fmt.Errorf("invalid pin: %d (min: %d, max: %d)", pin, MinPin, MaxPin)
		}
	}

	req := NewSetPinsRequest(id, pins)
	resp := &Response[struct{}]{}
	_, err := RunCommand(id, device, req, resp)
	return err
}

// GetDuty retrieves the current duty setting from the Picow device. Default Request ID will be used.
func GetDuty(device *models.Device) ([]uint8, error) {
	return GetDutyWithID(RequestIDDefault, device)
}

// GetDutyWithID retrieves the current duty setting from the Picow device.
func GetDutyWithID(id RequestID, device *models.Device) ([]uint8, error) {
	req := NewGetDutyRequest(id)
	resp := &Response[[]uint8]{}
	return RunCommand(id, device, req, resp)
}

// SetDuty sets the duty on the Picow device. Default request ID will be used.
func SetDuty(device *models.Device, duty ...uint8) error {
	return SetDutyWithID(RequestIDDefault, device, duty...)
}

// SetDutyWithID sets the duty on the Picow device.
func SetDutyWithID(id RequestID, device *models.Device, duty ...uint8) error {
	// Validate duty first
	for _, d := range duty {
		if d < MinDuty || d > MaxDuty {
			return fmt.Errorf("invalid duty: %d (min: %d, max: %d)", d, MinDuty, MaxDuty)
		}
	}

	req := NewSetDutyRequest(id, duty)
	resp := &Response[struct{}]{}
	_, err := RunCommand(id, device, req, resp)
	return err
}

// GetTemperature retrieves the current temperature from the Picow device. Default request ID will be used.
func GetTemperature(device *models.Device) (float64, error) {
	return GetTemperatureWithID(RequestIDDefault, device)
}

// GetTemperatureWithID retrieves the current temperature from the Picow device.
func GetTemperatureWithID(id RequestID, device *models.Device) (float64, error) {
	req := NewGetTemperatureRequest(id)
	resp := &Response[float64]{}
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
