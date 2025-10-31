package services

import (
	"reflect"
	"testing"

	"github.com/knackwurstking/picow-led/models"
)

func TestAddGroup(t *testing.T) {
	r := openDB(t, true)
	defer r.Close()

	_, err := r.Devices.Add(models.NewDevice("192.168.178.10:8888", "Test Device 1"))
	if err != nil {
		t.Fatalf("Failed to add device 1: %v", err)
	}

	r.Devices.Add(models.NewDevice("192.168.178.20:8888", "Test Device 2"))
	if err != nil {
		t.Fatalf("Failed to add device 2: %v", err)
	}

	r.Devices.Add(models.NewDevice("192.168.178.30:8888", "Test Device 3"))
	if err != nil {
		t.Fatalf("Failed to add device 3: %v", err)
	}

	// Create a new group - test failure, invalid devices
	group := models.NewGroup("Test Group", []models.DeviceID{1, 2, 3, 4})

	// Add the group to the database
	id, err := r.Groups.Add(group)
	if err != ErrInvalidDeviceID {
		t.Fatalf("Failed to add group: %#v", err)
	}

	// Create a new group
	group = models.NewGroup("Test Group", []models.DeviceID{1, 2, 3})

	// Add the group to the database
	id, err = r.Groups.Add(group)
	if err != nil {
		t.Fatalf("Failed to add group: %#v", err)
	}

	retrievedGroup, err := r.Groups.Get(id)
	if err != nil {
		t.Fatalf("Failed to retrieve group with ID %d", id)
	}
	if retrievedGroup.ID != id {
		t.Errorf("Expected group ID %d, got %d", id, retrievedGroup.ID)
	}
	if retrievedGroup.Name != group.Name {
		t.Errorf("Expected group %v, got %v", group, retrievedGroup)
	}
	if !reflect.DeepEqual(retrievedGroup.Devices, group.Devices) {
		t.Errorf("Expected group devices %#v, got %#v", group.Devices, retrievedGroup.Devices)
	}
}

func TestUpdateGroup(t *testing.T) {
	r := openDB(t, false)
	defer r.Close()

	// Update the group in the database - check fail
	group := models.NewGroup("Updated Group", []models.DeviceID{1, 2, 3, 4})
	group.ID = 1

	err := r.Groups.Update(group)
	if err != ErrInvalidDeviceID {
		t.Fatal("Expected error updating group with invalid devices")
	}

	// Update the group in the database
	group = models.NewGroup("Updated Group", []models.DeviceID{1, 2})
	group.ID = 1

	err = r.Groups.Update(group)
	if err != nil {
		t.Fatalf("Failed to update group: %v", err)
	}

	retrievedGroup, err := r.Groups.Get(1)
	if err != nil {
		t.Fatalf("Failed to retrieve group with ID 1")
	}
	if retrievedGroup.ID != 1 {
		t.Errorf("Expected group ID 1, got %d", retrievedGroup.ID)
	}
	if retrievedGroup.Name != group.Name {
		t.Errorf("Expected group %v, got %v", group, retrievedGroup)
	}
	if !reflect.DeepEqual(retrievedGroup.Devices, group.Devices) {
		t.Errorf("Expected group devices %#v, got %#v", group.Devices, retrievedGroup.Devices)
	}
}

func TestDeleteGroup(t *testing.T) {
	r := openDB(t, false)
	defer r.Close()

	// Delete the group from the database
	err := r.Groups.Delete(1)
	if err != nil {
		t.Fatal("Failed to delete group with ID 1")
	}

	// Retrieve the group from the database
	retrievedGroup, err := r.Groups.Get(1)
	if err == nil {
		t.Fatalf("Expected group to be deleted, got %v", retrievedGroup)
	}
}
