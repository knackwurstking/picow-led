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
	device1 := models.NewDevice("192.168.178.10:8888", "Test Device 1")
	device1.ID = models.DeviceID(1)

	id, err := r.Devices.Add(device1)
	if err != nil {
		t.Fatalf("Failed to add device: %v", err)
	}
	if id != device1.ID {
		t.Errorf("Expected ID 1, got %d", id)
	}

	dbDevice, err := r.Devices.Get(id)
	if err != nil {
		t.Fatalf("Failed to get device: %v", err)
	}
	if dbDevice.ID != id {
		t.Errorf("Expected ID %d, got %d", id, dbDevice.ID)
	}
	if dbDevice.Addr != device1.Addr {
		t.Errorf("Expected IP %v, got %s", device1.Addr, dbDevice.Addr)
	}
	if dbDevice.Name != device1.Name {
		t.Errorf("Expected Name %v, got %s", device1.Name, dbDevice.Name)
	}

	// Second device
	device2 := models.NewDevice("192.168.178.20:8888", "Test Device 2")
	device2.ID = models.DeviceID(2)

	id, err = r.Devices.Add(device2)
	if err != nil {
		t.Fatalf("Failed to add device: %v", err)
	}
	if id != device2.ID {
		t.Errorf("Expected ID %d, got %d", device2.ID, id)
	}
}

func TestAddDeviceSetup(t *testing.T) {
	r := openDB(t, false)
	defer r.Close()

	// Setup for the first device
	deviceSetup1 := models.NewDeviceSetup(1, []uint8{0, 1, 2, 3})

	id, err := r.DeviceSetups.Add(deviceSetup1)
	if err != nil {
		if !strings.Contains(err.Error(), "connect: no route to host") &&
			!strings.Contains(err.Error(), "timeout") {
			t.Fatalf("Failed to add pin: %v", err)
		}
	}
	if id != deviceSetup1.DeviceID {
		t.Errorf("Expected Device ID %d, got %d", deviceSetup1.DeviceID, id)
	}

	pins, err := r.DeviceSetups.Get(id)
	if err != nil {
		t.Fatalf("Failed to get pin: %v", err)
	}
	if pins.DeviceID != deviceSetup1.DeviceID {
		t.Errorf("Expected device ID %d, got %d", deviceSetup1.DeviceID, pins.DeviceID)
	}
	if !reflect.DeepEqual(pins.Pins, deviceSetup1.Pins) {
		t.Errorf("Expected pins %v, got %v", deviceSetup1.Pins, pins.Pins)
	}

	// Second setup for the first device
	// This should be not possible, adding 2 setups for one device
	deviceSetup2 := models.NewDeviceSetup(1, []uint8{4, 5, 6, 7})
	if _, err = r.DeviceSetups.Add(deviceSetup2); err == nil {
		t.Fatal("Expected an error while trying to add a second setup for device 1, got nil")
	}

	// Setup for the second device
	deviceSetup2 = models.NewDeviceSetup(2, []uint8{0, 1, 2, 3})
	if _, err = r.DeviceSetups.Add(deviceSetup2); err != nil {
		if !strings.Contains(err.Error(), "connect: no route to host") &&
			!strings.Contains(err.Error(), "timeout") {
			t.Fatalf("Failed to add device %d setup: %v", deviceSetup2.DeviceID, err)
		}
	}

	// Add a device control column for the first device
	deviceControl1 := models.NewDeviceControl(deviceSetup1.DeviceID, []uint8{255, 255, 255, 255})
	if _, err = r.DeviceControls.Add(deviceControl1); err != nil {
		t.Fatalf("Failed to add device %d control: %v", deviceControl1.DeviceID, err)
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
		t.Errorf("Expected not found error, the setup for device_id %d got not removed: %v", deviceID, err)
	}

	if _, err := r.DeviceControls.Get(deviceID); err != ErrNotFound {
		t.Errorf("Expected not found error, the device control for device_id %d got not removed: %v", deviceID, err)
	}
}

func TestDeviceSetupDelete(t *testing.T) {
	r := openDB(t, false)
	defer r.Close()

	deviceID := models.DeviceID(2)

	// Remove the setup for the second device
	if err := r.DeviceSetups.Delete(deviceID); err != nil {
		t.Fatalf("Failed to remove device setup: %v", err)
	}

	if _, err := r.DeviceSetups.Get(deviceID); err == nil {
		t.Errorf("Expected error, got nil, the pins with device_id %d got not removed", deviceID)
	}
}
