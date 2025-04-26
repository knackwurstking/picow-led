package api

import (
	"testing"
)

func TestMicroRequestError_InvalidServer(t *testing.T) {
	s := &Server{
		Addr: "1111:1111",
		Name: "Not existing Server",
	}
	r := &MicroRequest{
		ID:      MicroIDNoResponse,
		Group:   MicroGroupConfig,
		Command: "not-existing-command",
	}
	r.Send(s)
	if s.Error == "" {
		t.Errorf("Error expected, because the server address is invalid, but go nothing")
	} else {
		t.Log(s.Error)
	}
}
