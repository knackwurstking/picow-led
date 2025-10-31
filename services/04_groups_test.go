package services

import (
	"reflect"
	"testing"

	"github.com/knackwurstking/picow-led/models"
)

func TestAddGroup(t *testing.T) {
	r := openDB(t, true)
	defer r.Close()

	testDevice1Address := models.Addr("192.168.178.10:8888")
	testDevice1Name := "Test Device 1"
	testDevice1 := models.NewDevice(testDevice1Address, testDevice1Name)

	if _, err := r.Devices.Add(testDevice1); err != nil {
		t.Fatalf("Failed to add device 1: %v", err)
	}

	testDevice2Address := models.Addr("192.168.178.20:8888")
	testDevice2Name := "Test Device 2"
	testDevice2 := models.NewDevice(testDevice2Address, testDevice2Name)

	if _, err := r.Devices.Add(testDevice2); err != nil {
		t.Fatalf("Failed to add device 2: %v", err)
	}

	testDevice3Address := models.Addr("192.168.178.30:8888")
	testDevice3Name := "Test Device 3"
	testDevice3 := models.NewDevice(testDevice3Address, testDevice3Name)

	if _, err := r.Devices.Add(testDevice3); err != nil {
		t.Fatalf("Failed to add device 3: %v", err)
	}

	// Create a new group - test failure, invalid devices
	testGroup1Name := "Test Group 1"
	testGroup1Devices := []models.DeviceID{1, 2, 3, 4}
	testGroup1 := models.NewGroup(testGroup1Name, testGroup1Devices)

	// Add the group to the database
	if _, err := r.Groups.Add(testGroup1); err != ErrInvalidDeviceID {
		t.Fatalf("Failed to add group: %#v", err)
	}

	// Create a new group
	testGroup2Name := "Test Group 2"
	testGroup2Devices := []models.DeviceID{1, 2, 3}
	testGroup2 := models.NewGroup(testGroup2Name, testGroup2Devices)

	// Add the group to the database
	if id, err := r.Groups.Add(testGroup2); err != nil {
		t.Fatalf("Failed to add group: %#v", err)
	} else {
		retrievedGroup, err := r.Groups.Get(id)
		if err != nil {
			t.Fatalf("Failed to retrieve group with ID %d", id)
		}
		if retrievedGroup.ID != id {
			t.Errorf("Expected group ID %d, got %d", id, retrievedGroup.ID)
		}
		if retrievedGroup.Name != testGroup2.Name {
			t.Errorf("Expected group %v, got %v", testGroup2, retrievedGroup)
		}
		if !reflect.DeepEqual(retrievedGroup.Devices, testGroup2.Devices) {
			t.Errorf("Expected group devices %#v, got %#v", testGroup2.Devices, retrievedGroup.Devices)
		}
	}
}

// TODO: Continue refactoring (clean up) here
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
