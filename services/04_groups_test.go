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
	testDevice1, err := models.NewDevice(testDevice1Address, testDevice1Name)
	if err != nil {
		t.Fatalf("Failed to create device 1: %v", err)
	}

	if _, err := r.Devices.Add(testDevice1); err != nil {
		t.Fatalf("Failed to add device 1: %v", err)
	}

	testDevice2Address := models.Addr("192.168.178.20:8888")
	testDevice2Name := "Test Device 2"
	testDevice2, err := models.NewDevice(testDevice2Address, testDevice2Name)
	if err != nil {
		t.Fatalf("Failed to create device 2: %v", err)
	}

	if _, err := r.Devices.Add(testDevice2); err != nil {
		t.Fatalf("Failed to add device 2: %v", err)
	}

	testDevice3Address := models.Addr("192.168.178.30:8888")
	testDevice3Name := "Test Device 3"
	testDevice3, err := models.NewDevice(testDevice3Address, testDevice3Name)
	if err != nil {
		t.Fatalf("Failed to create device 3: %v", err)
	}

	if _, err := r.Devices.Add(testDevice3); err != nil {
		t.Fatalf("Failed to add device 3: %v", err)
	}

	// Create a new group - test failure, invalid devices
	testGroup1Name := "Test Group 1"
	testGroup1Devices := []models.DeviceID{1, 2, 3, 4}
	testGroup1 := models.NewGroup(testGroup1Name, testGroup1Devices)

	// Add the group to the database
	if _, err := r.Groups.Add(testGroup1); err != ErrInvalidGroupDeviceID {
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

func TestUpdateGroup(t *testing.T) {
	r := openDB(t, false)
	defer r.Close()

	// Update the group in the database - check fail
	group1Name := "Updated Group"
	group1Devices := []models.DeviceID{1, 2, 3, 4}
	group1 := models.NewGroup(group1Name, group1Devices)
	group1.ID = 1

	err := r.Groups.Update(group1)
	if err != ErrInvalidGroupDeviceID {
		t.Fatal("Expected error updating group with invalid devices")
	}

	// Update the group in the database
	group1.Devices = []models.DeviceID{1, 2}

	if err = r.Groups.Update(group1); err != nil {
		t.Fatalf("Failed to update group: %v", err)
	}

	retrievedGroup, err := r.Groups.Get(group1.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve group with ID 1")
	}
	if retrievedGroup.ID != group1.ID {
		t.Errorf("Expected group ID %d, got %d", group1.ID, retrievedGroup.ID)
	}
	if retrievedGroup.Name != group1.Name {
		t.Errorf("Expected group %v, got %v", group1.Name, retrievedGroup.Name)
	}
	if !reflect.DeepEqual(retrievedGroup.Devices, group1.Devices) {
		t.Errorf("Expected group devices %#v, got %#v", group1.Devices, retrievedGroup.Devices)
	}
}

func TestDeleteGroup(t *testing.T) {
	r := openDB(t, false)
	defer r.Close()

	groupID := models.GroupID(1)

	// Delete the group from the database
	if err := r.Groups.Delete(groupID); err != nil {
		t.Fatalf("Failed to delete group with ID %d: %v", groupID, err)
	}

	// Retrieve the group from the database
	if retrievedGroup, err := r.Groups.Get(groupID); err == nil {
		t.Fatalf("Expected group to be deleted, got %v", retrievedGroup)
	}
}
