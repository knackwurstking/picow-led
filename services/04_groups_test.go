package services

import (
	"reflect"
	"testing"

	"github.com/knackwurstking/picow-led/models"
)

func TestAddGroup(t *testing.T) {
	r := openDB(t, true)
	defer r.Close()

	// Create a new group
	group := &models.Group{
		Name: "Test Group",
		Setup: models.GroupSetup{
			{DeviceID: 1, ColorID: 1},
			{DeviceID: 2, ColorID: 2},
			{DeviceID: 3, ColorID: 3},
		},
	}

	// Add the group to the database
	id, err := r.Groups.Add(group)
	if err != nil {
		t.Fatalf("Failed to add group: %v", err)
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
	if !reflect.DeepEqual(retrievedGroup.Setup, group.Setup) {
		t.Errorf("Expected group setup %#v, got %#v", group.Setup, retrievedGroup.Setup)
	}
}

func TestUpdateGroup(t *testing.T) {
	r := openDB(t, false)
	defer r.Close()

	// Update the group in the database
	group := &models.Group{
		ID:   1,
		Name: "Updated Group",
		Setup: models.GroupSetup{
			{DeviceID: 1, ColorID: 1},
			{DeviceID: 2, ColorID: 2},
			{DeviceID: 3, ColorID: 4},
		},
	}

	err := r.Groups.Update(group)
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
	if !reflect.DeepEqual(retrievedGroup.Setup, group.Setup) {
		t.Errorf("Expected group setup %#v, got %#v", group.Setup, retrievedGroup.Setup)
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
