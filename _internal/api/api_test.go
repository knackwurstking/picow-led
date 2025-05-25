package api

import (
	"picow-led/internal/types"
	"testing"
)

func TestMicroRequestError_InvalidCommand(t *testing.T) {
	apiConfig, _ := GetConfig(".api.yaml")
	devices := []*types.Device{}
	for _, d := range apiConfig.Devices {
		devices = append(devices, d)
	}

	r := &types.MicroRequest{
		ID:      types.MicroIDDefault,
		Type:    types.MicroTypeGET,
		Group:   types.MicroGroupConfig,
		Command: "not-existing-command",
	}

	for _, d := range devices {
		respData, err := r.Send(d)
		if err != nil {
			continue
		}

		_, err = types.ParseMicroResponse[any](respData)
		if err == nil {
			t.Errorf("Error expected, because the command does not exists, but go nothing. server=%+v", d)
		} else {
			t.Log(err)
		}
	}
}

func TestMicroRequestError_InvalidGroup(t *testing.T) {
	apiConfig, _ := GetConfig(".api.yaml")
	devices := []*types.Device{}
	for _, d := range apiConfig.Devices {
		devices = append(devices, d)
	}

	r := &types.MicroRequest{
		ID:      types.MicroIDDefault,
		Type:    types.MicroTypeGET,
		Group:   "wrong-group",
		Command: "not-existing-command",
	}

	for _, d := range devices {
		respData, err := r.Send(d)
		if err != nil {
			continue
		}

		_, err = types.ParseMicroResponse[any](respData)
		if err == nil {
			t.Errorf("Error expected, because the command does not exists, but go nothing. server=%+v", d)
		} else {
			t.Log(err)
		}
	}
}

func TestMicroRequestError_InvalidType(t *testing.T) {
	apiConfig, _ := GetConfig(".api.yaml")
	devices := []*types.Device{}
	for _, d := range apiConfig.Devices {
		devices = append(devices, d)
	}

	r := &types.MicroRequest{
		ID:      types.MicroIDDefault,
		Type:    "wrong-type",
		Group:   "wrong-group",
		Command: "not-existing-command",
	}

	for _, d := range devices {
		respData, err := r.Send(d)
		if err != nil {
			continue
		}

		_, err = types.ParseMicroResponse[any](respData)
		if err == nil {
			t.Errorf("Error expected, because the command does not exists, but go nothing. server=%+v", d)
		} else {
			t.Log(err)
		}
	}
}
