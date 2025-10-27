// Testing device and setup creation and removal
package services

import (
	"testing"

	"github.com/knackwurstking/picow-led/models"

	_ "github.com/mattn/go-sqlite3"
)

func TestAddDevice(t *testing.T) {
	db := openDB(t, false)
	defer db.Close()

	// First device
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

	// Second device
	id, err = registry.Devices.Add(models.NewDevice("192.168.178.20", "Test Device 2"))
	if err != nil {
		t.Fatalf("Failed to add device: %v", err)
	}
	if id != 2 {
		t.Errorf("Expected ID 2, got %d", id)
	}
}

func TestAddDeviceSetup(t *testing.T) {
	db := openDB(t, false)
	defer db.Close()

	// Setup for the first device
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

	// Second setup for the first device
	// This should be not possible, adding 2 setups for one device
	_, err = registry.DeviceSetups.Add(models.NewDeviceSetup(1, []uint8{4, 5, 6, 7}))
	if err == nil {
		t.Fatal("Expected an error while trying to add a second setup for device 1, got nil")
	}

	// Setup for the second device
	_, err = registry.DeviceSetups.Add(models.NewDeviceSetup(2, []uint8{0, 1, 2, 3}))
	if err != nil {
		t.Fatalf("Failed to add device 2 setup: %v", err)
	}
}

func TestRemoveDevice(t *testing.T) {
	db := openDB(t, false)
	defer db.Close()

	// Remove device 1 for checking if the setup got removed too
	if err := registry.Devices.Delete(1); err != nil {
		t.Fatalf("Failed to remove device: %v", err)
	}

	if _, err := registry.DeviceSetups.Get(1); err == nil {
		t.Errorf("Expected error, got nil, the pins with device_id 1 got not removed")
	}
}

func TestDeviceSetupDelete(t *testing.T) {
	db := openDB(t, false)
	defer db.Close()

	// Remove the setup for the second device
	if err := registry.DeviceSetups.Delete(2); err != nil {
		t.Fatalf("Failed to remove device setup: %v", err)
	}

	if _, err := registry.DeviceSetups.Get(2); err == nil {
		t.Error("Expected error, got nil, the pins with device_id 2 got not removed")
	}
}
