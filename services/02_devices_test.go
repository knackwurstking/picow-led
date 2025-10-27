package services

import (
	"testing"

	"github.com/knackwurstking/picow-led/models"

	_ "github.com/mattn/go-sqlite3"
)

func TestAddDevice(t *testing.T) {
	db := openDB(t)
	defer db.Close()

	id, err := registry.Devices.Add(models.NewDevice("192.168.178.10", "Test Device 1"))
	if err != nil {
		t.Fatalf("Failed to add device: %v", err)
	}
	if id != 1 {
		t.Errorf("Expected ID 1, got %d", id)
	}

	device, err := registry.Devices.Get(id)
	if err != nil {
		t.Fatalf("Failed to get device: %v", err)
	}
	if device.ID != id {
		t.Errorf("Expected ID %d, got %d", id, device.ID)
	}
	if device.Addr != "192.168.178.10" {
		t.Errorf("Expected IP 192.168.178.10, got %s", device.Addr)
	}
	if device.Name != "Test Device 1" {
		t.Errorf("Expected Name Test Device 1, got %s", device.Name)
	}
}

func TestAddDeviceSetup(t *testing.T) {
	db := openDB(t)
	defer db.Close()

	id, err := registry.DeviceSetups.Add(models.NewDeviceSetup(1, []uint8{0, 1, 2, 3}))
	if err != nil {
		t.Fatalf("Failed to add pin: %v", err)
	}
	if id != 1 {
		t.Errorf("Expected ID 1, got %d", id)
	}

	pins, err := registry.DeviceSetups.Get(id)
	if err != nil {
		t.Fatalf("Failed to get pin: %v", err)
	}
	if pins.DeviceID != 1 {
		t.Errorf("Expected device ID 1, got %d", pins.DeviceID)
	}
	if len(pins.Pins) != 4 {
		t.Errorf("Expected 4 pins, got %d", len(pins.Pins))
	}
}

func TestRemoveDevice(t *testing.T) {
	db := openDB(t)
	defer db.Close()

	if err := registry.Devices.Delete(1); err != nil {
		t.Fatalf("Failed to remove device: %v", err)
	}

	if _, err := registry.DeviceSetups.Get(1); err == nil {
		t.Errorf("Expected error, got nil, the pins with device_id 1 got not removed")
	}
}
