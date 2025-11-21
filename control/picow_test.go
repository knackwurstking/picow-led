package control

import (
	"bytes"
	"encoding/json"
	"net"
	"testing"

	"github.com/knackwurstking/picow-led/models"
)

func TestEndByte(t *testing.T) {
	device, err := models.NewDevice("192.168.178.10:8888", "Test Device 1")
	if err != nil {
		t.Fatalf("Failed to create device: %v", err)
	}

	picow := NewPicoW(device)

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

	if result := picow.EndByte(data); !bytes.Equal(result[len(result)-1:], []byte("\n")) {
		t.Errorf("Invalid last byte, got %#v != %#v", result[len(result)-1], []byte("\n")[0])
	}
}

func TestPicoWReadFromConnUntilEndByte(t *testing.T) {
	device, err := models.NewDevice("192.168.178.10:8888", "Test Device 1")
	if err != nil {
		t.Fatalf("Failed to create device: %v", err)
	}

	picow := NewPicoW(device)

	server, client := net.Pipe()
	picow.Conn = client
	defer client.Close()
	defer server.Close()

	go func() {
		defer server.Close()

		_, err := server.Write([]byte("test\n"))
		if err != nil {
			t.Errorf("Failed to write to server: %v", err)
		}
	}()

	data, err := picow.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read from connection: %v", err)
	}

	if !bytes.Equal(data, []byte("test")) {
		t.Errorf("Invalid read result, got %#v != %#v", data, []byte("test"))
	}
}
