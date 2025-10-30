package control

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/knackwurstking/picow-led/models"
)

func TestEndByte(t *testing.T) {
	device := NewDevice(
		models.NewDevice("192.168.178.10:8888", "Test Device 1"),
		models.NewDeviceSetup(1, []uint8{1, 2, 3, 4}),
	)

	r := NewRequest(
		RequestIDDefault,
		TypeGet,
		"led",
		"color",
	)

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	dataOriginal := make([]byte, len(data))
	copy(dataOriginal, data)

	if result := device.EndByte(data); !bytes.Equal(result[len(result)-1:], []byte("\n")) {
		t.Errorf("Invalid last byte, got %#v != %#v", result[len(result)-1], []byte("\n")[0])
	}
}
