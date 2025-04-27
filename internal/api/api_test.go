package api

import (
	"testing"
)

func TestMicroRequestError_InvalidCommand(t *testing.T) {
	apiOptions, _ := GetApiOptions(".api.yaml")

	r := &MicroRequest{
		ID:      MicroIDDefault,
		Type:    MicroTypeGET,
		Group:   MicroGroupConfig,
		Command: "not-existing-command",
	}

	for _, s := range apiOptions.Servers {
		respData, err := r.Send(s)
		if err != nil {
			continue
		}

		_, err = ParseMicroResponse[any](respData)
		if err != nil {
			s.Error = err.Error()
		}

		if s.Error == "" {
			t.Errorf("Error expected, because the command does not exists, but go nothing. server=%+v", s)
		} else {
			t.Log(s.Error)
		}
	}
}

func TestMicroRequestError_InvalidGroup(t *testing.T) {
	apiOptions, _ := GetApiOptions(".api.yaml")

	r := &MicroRequest{
		ID:      MicroIDDefault,
		Type:    MicroTypeGET,
		Group:   "wrong-group",
		Command: "not-existing-command",
	}

	for _, s := range apiOptions.Servers {
		respData, err := r.Send(s)
		if err != nil {
			continue
		}

		_, err = ParseMicroResponse[any](respData)
		if err != nil {
			s.Error = err.Error()
		}

		if s.Error == "" {
			t.Errorf("Error expected, because the command does not exists, but go nothing. server=%+v", s)
		} else {
			t.Log(s.Error)
		}
	}
}

func TestMicroRequestError_InvalidType(t *testing.T) {
	apiOptions, _ := GetApiOptions(".api.yaml")

	r := &MicroRequest{
		ID:      MicroIDDefault,
		Type:    "wrong-type",
		Group:   "wrong-group",
		Command: "not-existing-command",
	}

	for _, s := range apiOptions.Servers {
		respData, err := r.Send(s)
		if err != nil {
			continue
		}

		_, err = ParseMicroResponse[any](respData)
		if err != nil {
			s.Error = err.Error()
		}

		if s.Error == "" {
			t.Errorf("Error expected, because the command does not exists, but go nothing. server=%+v", s)
		} else {
			t.Log(s.Error)
		}
	}
}
