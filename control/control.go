package control

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/knackwurstking/picow-led/models"
)

const (
	MinDuty uint8 = 0
	MaxDuty uint8 = 255

	MinPin uint8 = 0
	MaxPin uint8 = 15
)

var (
	ErrNotConnected = errors.New("not connected")
	ErrNoData       = errors.New("no data")
)

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

func GetPins(id RequestID, device *models.ResolvedDevice) ([]uint8, error) {
	req := NewGetPinsRequest(id)
	resp := &Response[[]uint8]{}
	return RunCommand(id, device, req, resp)
}

func SetPins(id RequestID, device *models.ResolvedDevice, pins ...uint8) error {
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

func GetColor(id RequestID, device *models.ResolvedDevice) ([]uint8, error) {
	req := NewGetColorRequest(id)
	resp := &Response[[]uint8]{}
	return RunCommand(id, device, req, resp)
}

func SetColor(id RequestID, device *models.ResolvedDevice, duty ...uint8) error {
	// Validate duty first
	for _, d := range duty {
		if d < MinDuty || d > MaxDuty {
			return fmt.Errorf("invalid duty: %d (min: %d, max: %d)", d, MinDuty, MaxDuty)
		}
	}

	req := NewSetColorRequest(id, duty)
	resp := &Response[struct{}]{}
	_, err := RunCommand(id, device, req, resp)
	return err
}

func GetTemperature(id RequestID, device *models.ResolvedDevice) (float32, error) {
	req := NewGetTemperatureRequest(id)
	resp := &Response[float32]{}
	return RunCommand(id, device, req, resp)
}

func GetDiskUsage(id RequestID, device *models.ResolvedDevice) (*DiskUsage, error) {
	req := NewGetDiskUsageRequest(id)
	resp := &Response[*DiskUsage]{}
	return RunCommand(id, device, req, resp)
}

func GetVersion(id RequestID, device *models.ResolvedDevice) (string, error) {
	req := NewGetVersionRequest(id)
	resp := &Response[string]{}
	return RunCommand(id, device, req, resp)
}
