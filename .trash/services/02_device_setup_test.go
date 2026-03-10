// Testing device and setup creation and removal
package services

import (
	"testing"

	"github.com/knackwurstking/picow-led/models"

	_ "github.com/mattn/go-sqlite3"
)

func TestAddDevice(t *testing.T) {
	r := openDB(t, true)
	defer r.Close()

	// First device
	device1, err := models.NewDevice("192.168.178.10:8888", "Test Device 1")
	if err != nil {
		t.Fatalf("create device: %v", err)
	}
	device1.ID = models.DeviceID(1)

	id, err := r.Devices.Add(device1)
	if err != nil {
		t.Fatalf("add device: %v", err)
	}
	if id != device1.ID {
		t.Errorf("Expected ID 1, got %d", id)
	}

	dbDevice, err := r.Devices.Get(id)
	if err != nil {
		t.Fatalf("get device: %v", err)
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
	device2, err := models.NewDevice("192.168.178.20:8888", "Test Device 2")
	if err != nil {
		t.Fatalf("create device: %v", err)
	}
	device2.ID = models.DeviceID(2)

	id, err = r.Devices.Add(device2)
	if err != nil {
		t.Fatalf("add device: %v", err)
	}
	if id != device2.ID {
		t.Errorf("Expected ID %d, got %d", device2.ID, id)
	}
}

func TestRemoveDevice(t *testing.T) {
	r := openDB(t, false)
	defer r.Close()

	deviceID := models.DeviceID(1)

	// Remove device 1 for checking if the setup got removed too
	if err := r.Devices.Delete(deviceID); err != nil {
		t.Fatalf("remove device: %v", err)
	}

	if _, err := r.DeviceControls.Get(deviceID); !IsNotFoundError(err) {
		t.Errorf("Expected not found error, the device control for device_id %d got not removed: %v", deviceID, err)
	}
}
