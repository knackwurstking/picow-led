package services

import (
	"reflect"
	"testing"

	"github.com/knackwurstking/picow-led/models"
)

func TestAddColor(t *testing.T) {
	r := openDB(t, true)
	defer r.Close()

	color := &models.Color{
		Name: "Test Color",
		Duty: models.Duty{255, 255, 255, 255},
	}

	id, err := r.Colors.Add(color)
	if err != nil {
		t.Fatalf("Failed to add color: %v", err)
	}
	if id != 1 {
		t.Errorf("Expected color ID 1, got %d", id)
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
