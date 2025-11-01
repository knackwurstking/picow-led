package services

import (
	"reflect"
	"testing"

	"github.com/knackwurstking/picow-led/models"
)

func TestAddColor(t *testing.T) {
	r := openDB(t, true)
	defer r.Close()

	color := models.NewColor("Test Color", []uint8{255, 255, 255, 255})
	color.ID = models.ColorID(1)

	id, err := r.Colors.Add(color)
	if err != nil {
		t.Fatalf("Failed to add color: %v", err)
	}
	if id != color.ID {
		t.Errorf("Expected color ID %d, got %d", color.ID, id)
	}

	// Verify that the color was added to the database
	dbColor, err := r.Colors.Get(id)
	if err != nil {
		t.Fatalf("Failed to get color: %v", err)
	}
	if dbColor.ID != id {
		t.Errorf("Expected color ID %d, got %d", id, dbColor.ID)
	}
	if dbColor.Name != color.Name {
		t.Errorf("Expected color name %s, got %s", color.Name, dbColor.Name)
	}
	if !reflect.DeepEqual(dbColor.Duty, color.Duty) {
		t.Errorf("Expected color duty %v, got %v", color.Duty, dbColor.Duty)
	}
}

func TestUpdateColor(t *testing.T) {
	r := openDB(t, false)
	defer r.Close()

	color := models.NewColor("Updated Color", []uint8{128, 128, 128, 128})
	color.ID = models.ColorID(1)

	err := r.Colors.Update(color)
	if err != nil {
		t.Fatalf("Failed to update color: %v", err)
	}

	// Verify that the color was updated in the database
	dbColor, err := r.Colors.Get(color.ID)
	if err != nil {
		t.Fatalf("Failed to get color: %v", err)
	}
	if dbColor.ID != color.ID {
		t.Errorf("Expected color ID %d, got %d", color.ID, dbColor.ID)
	}
	if dbColor.Name != color.Name {
		t.Errorf("Expected color name %s, got %s", color.Name, dbColor.Name)
	}
	if !reflect.DeepEqual(dbColor.Duty, color.Duty) {
		t.Errorf("Expected color duty %v, got %v", color.Duty, dbColor.Duty)
	}
}
