package services

import (
	"log/slog"

	"github.com/knackwurstking/picow-led/models"
)

type Pins struct {
	registry *Registry
}

func NewPins(registry *Registry) *Pins {
	return &Pins{
		registry: registry,
	}
}

func (p *Pins) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS pins (
		device_id INTEGER UNIQUE NOT NULL,
		name TEXT NOT NULL,
		pins TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := p.registry.db.Exec(query)
	return err
}

func (p *Pins) Get(id models.DeviceID) (*models.Pins, error) {
	slog.Debug("Get pins from database", "table", "pins", "device_id", id)

	query := `SELECT * FROM pins WHERE device_id = ?`
	pins, err := ScanPins(p.registry.db.QueryRow(query, id))
	if err != nil {
		return nil, err
	}

	return pins, nil
}

func (p *Pins) List() ([]*models.Pins, error) {
	slog.Debug("List pins from database", "table", "pins")

	query := `SELECT * FROM pins`
	rows, err := p.registry.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pins []*models.Pins
	for rows.Next() {
		pin, err := ScanPins(rows)
		if err != nil {
			return nil, err
		}
		pins = append(pins, pin)
	}

	return pins, nil
}

func (p *Pins) Update(pins *models.Pins) error {
	slog.Debug("Update pins in database", "table", "pins", "device_id", pins.DeviceID)

	query := `UPDATE pins SET name = ?, pins = ? WHERE device_id = ?`
	_, err := p.registry.db.Exec(query, pins.Name, pins.Pins, pins.DeviceID)
	return err
}

func (p *Pins) Delete(id models.DeviceID) error {
	slog.Debug("Delete pins from database", "table", "pins", "device_id", id)

	query := `DELETE FROM pins WHERE device_id = ?`
	_, err := p.registry.db.Exec(query, id)
	return err
}

func ScanPins(r Scannable) (*models.Pins, error) {
	pins := &models.Pins{}

	if err := r.Scan(&pins.DeviceID, &pins.Name, &pins.Pins, &pins.CreatedAt); err != nil {
		return nil, err
	}

	return pins, nil
}
