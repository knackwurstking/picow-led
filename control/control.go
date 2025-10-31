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

func GetPins(id RequestID, device *models.ResolvedDevice) ([]uint8, error) {
	req := NewGetPinsRequest(id)
	picow := NewPicoW(device)

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}
	n, err := picow.Write(data)
	if err != nil {
		return nil, err
	}
	if n != len(data) {
		return nil, ErrNoData
	}

	if id == RequestIDNoResponse {
		return nil, nil
	}

	respData, err := picow.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	var resp *GetPinsResponse
	if err = json.Unmarshal(respData, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("picow error: %s", resp.Error)
	}

	return resp.Data, nil
}

func SetPins(id RequestID, device *models.ResolvedDevice, pins ...uint8) error {
	// Validate pins first
	for _, pin := range pins {
		if pin < MinPin || pin > MaxPin {
			return fmt.Errorf("invalid pin: %d (min: %d, max: %d)", pin, MinPin, MaxPin)
		}
	}

	req := NewSetPinsRequest(id, pins)
	picow := NewPicoW(device)

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}
	n, err := picow.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return ErrNoData
	}

	if id == RequestIDNoResponse {
		return nil
	}

	respData, err := picow.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}
	var resp *SetPinsResponse
	if err = json.Unmarshal(respData, &resp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}
	if resp.Error != "" {
		return fmt.Errorf("picow error: %s", resp.Error)
	}

	return nil
}

func GetColor(id RequestID, device *models.ResolvedDevice) ([]uint8, error) {
	req := NewGetColorRequest(id)
	picow := NewPicoW(device)

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}
	n, err := picow.Write(data)
	if err != nil {
		return nil, err
	}
	if n != len(data) {
		return nil, ErrNoData
	}

	if id == RequestIDNoResponse {
		return nil, nil
	}

	respData, err := picow.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	var resp *GetColorResponse
	if err = json.Unmarshal(respData, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("picow error: %s", resp.Error)
	}

	return resp.Data, nil
}

func SetColor(id RequestID, device *models.ResolvedDevice, duty ...uint8) error {
	// Validate duty first
	for _, d := range duty {
		if d < MinDuty || d > MaxDuty {
			return fmt.Errorf("invalid duty: %d (min: %d, max: %d)", d, MinDuty, MaxDuty)
		}
	}

	req := NewSetColorRequest(id, duty)
	picow := NewPicoW(device)

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}
	n, err := picow.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return ErrNoData
	}

	if id == RequestIDNoResponse {
		return nil
	}

	respData, err := picow.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}
	var resp *SetColorResponse
	if err = json.Unmarshal(respData, &resp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}
	if resp.Error != "" {
		return fmt.Errorf("picow error: %s", resp.Error)
	}

	return nil
}
