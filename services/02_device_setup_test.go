// Testing device and setup creation and removal
package services

import (
	"reflect"
	"strings"
	"testing"

	"github.com/knackwurstking/picow-led/models"

	_ "github.com/mattn/go-sqlite3"
)

func TestAddDevice(t *testing.T) {
	r := openDB(t, true)
	defer r.Close()

	// First device
	id, err := r.Devices.Add(models.NewDevice("192.168.178.10:8888", "Test Device 1"))
	if err != nil {
		t.Fatalf("Failed to add device: %v", err)
	}
	if id != 1 {
		t.Errorf("Expected ID 1, got %d", id)
	}

	device, err := r.Devices.Get(id)
	if err != nil {
		t.Fatalf("Failed to get device: %v", err)
	}
	if device.ID != id {
		t.Errorf("Expected ID %d, got %d", id, device.ID)
	}
	if device.Addr != "192.168.178.10:8888" {
		t.Errorf("Expected IP 192.168.178.10:8888, got %s", device.Addr)
	}
	if device.Name != "Test Device 1" {
		t.Errorf("Expected Name Test Device 1, got %s", device.Name)
	}

	// Second device
	id, err = r.Devices.Add(models.NewDevice("192.168.178.20:8888", "Test Device 2"))
	if err != nil {
		t.Fatalf("Failed to add device: %v", err)
	}
	if id != 2 {
		t.Errorf("Expected ID 2, got %d", id)
	}
}

func TestAddDeviceSetup(t *testing.T) {
	r := openDB(t, false)
	defer r.Close()

	// Setup for the first device
	deviceSetup := models.NewDeviceSetup(1, []uint8{0, 1, 2, 3})

	id, err := r.DeviceSetups.Add(deviceSetup)
	if err != nil {
		if !strings.Contains(err.Error(), "connect: no route to host") && !strings.Contains(err.Error(), "timeout") {
			t.Fatalf("Failed to add pin: %v", err)
		}
	}
	if id != 1 {
		t.Errorf("Expected ID 1, got %d", id)
	}

	pins, err := r.DeviceSetups.Get(id)
	if err != nil {
		t.Fatalf("Failed to get pin: %v", err)
	}
	if pins.DeviceID != 1 {
		t.Errorf("Expected device ID 1, got %d", pins.DeviceID)
	}
	if !reflect.DeepEqual(pins.Pins, deviceSetup.Pins) {
		t.Errorf("Expected pins %v, got %v", deviceSetup.Pins, pins.Pins)
	}

	// Second setup for the first device
	// This should be not possible, adding 2 setups for one device
	deviceSetup = models.NewDeviceSetup(1, []uint8{4, 5, 6, 7})
	_, err = r.DeviceSetups.Add(deviceSetup)
	if err == nil {
		t.Fatal("Expected an error while trying to add a second setup for device 1, got nil")
	}

	// Setup for the second device
	_, err = r.DeviceSetups.Add(models.NewDeviceSetup(2, []uint8{0, 1, 2, 3}))
	if err != nil {
		if !strings.Contains(err.Error(), "connect: no route to host") && !strings.Contains(err.Error(), "timeout") {
			t.Fatalf("Failed to add device 2 setup: %v", err)
		}
	}

	// Add a device control column for the first device
	_, err = r.DeviceControls.Add(models.NewDeviceControl(1, []uint8{255, 255, 255, 255}))
	if err != nil {
		t.Fatalf("Failed to add device 1 control: %v", err)
	}
}

func TestRemoveDevice(t *testing.T) {
	r := openDB(t, false)
	defer r.Close()

	deviceID := models.DeviceID(1)

	// Remove device 1 for checking if the setup got removed too
	if err := r.Devices.Delete(deviceID); err != nil {
		t.Fatalf("Failed to remove device: %v", err)
	}

	if _, err := r.DeviceSetups.Get(deviceID); err != ErrNotFound {
		t.Errorf("Expected not found error, the setup for device_id 1 got not removed: %v", err)
	}

	if _, err := r.DeviceControls.Get(deviceID); err != ErrNotFound {
		t.Errorf("Expected not found error, the device control for device_id 1 got not removed: %v", err)
	}
}

func TestDeviceSetupDelete(t *testing.T) {
	r := openDB(t, false)
	defer r.Close()

	// Remove the setup for the second device
	if err := r.DeviceSetups.Delete(2); err != nil {
		t.Fatalf("Failed to remove device setup: %v", err)
	}

	if _, err := r.DeviceSetups.Get(2); err == nil {
		t.Error("Expected error, got nil, the pins with device_id 2 got not removed")
	}
}
