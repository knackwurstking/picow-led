package services

import (
	"database/sql"
	"log/slog"

	"github.com/knackwurstking/picow-led/models"
)

type DeviceSetups struct {
	registry *Registry
}

func NewDeviceSetups(registry *Registry) *DeviceSetups {
	return &DeviceSetups{
		registry: registry,
	}
}

func (p *DeviceSetups) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS device_setups (
		device_id INTEGER UNIQUE NOT NULL,
		pins TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := p.registry.db.Exec(query)
	return err
}

func (p *DeviceSetups) Get(id models.DeviceID) (*models.DeviceSetup, error) {
	slog.Debug("Get device setup from database", "table", "device_setups", "device_id", id)

	query := `SELECT * FROM device_setups WHERE device_id = ?`
	devicesSetup, err := ScanDevicesSetup(p.registry.db.QueryRow(query, id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return devicesSetup, nil
}

func (p *DeviceSetups) List() ([]*models.DeviceSetup, error) {
	slog.Debug("List device setup from database", "table", "device_setups")

	query := `SELECT * FROM device_setups`
	rows, err := p.registry.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deviceSetups []*models.DeviceSetup
	for rows.Next() {
		devicesSetup, err := ScanDevicesSetup(rows)
		if err != nil {
			return nil, err
		}
		deviceSetups = append(deviceSetups, devicesSetup)
	}

	return deviceSetups, nil
}

func (p *DeviceSetups) Add(deviceSetup *models.DeviceSetup) (models.DeviceID, error) {
	slog.Debug("Add device setup to database", "table", "device_setups", "device_id", deviceSetup.DeviceID)

	if !deviceSetup.Validate() {
		return 0, ErrInvalidDeviceSetup
	}

	query := `INSERT INTO device_setups (device_id, pins) VALUES (?, ?)`
	result, err := p.registry.db.Exec(query, deviceSetup.DeviceID, deviceSetup.Pins)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return models.DeviceID(id), nil
}

func (p *DeviceSetups) Update(deviceSetup *models.DeviceSetup) error {
	slog.Debug("Update device setup in database", "table", "device_setups", "device_id", deviceSetup.DeviceID)

	if !deviceSetup.Validate() {
		return ErrInvalidDeviceSetup
	}

	query := `UPDATE device_setups SET pins = ? WHERE device_id = ?`
	_, err := p.registry.db.Exec(query, deviceSetup.Pins, deviceSetup.DeviceID)
	return err
}

func (p *DeviceSetups) Delete(id models.DeviceID) error {
	slog.Debug("Delete device setup from database", "table", "device_setups", "device_id", id)

	query := `DELETE FROM device_setups WHERE device_id = ?`
	_, err := p.registry.db.Exec(query, id)
	return err
}

func ScanDevicesSetup(r Scannable) (*models.DeviceSetup, error) {
	deviceSetup := &models.DeviceSetup{}

	if err := r.Scan(&deviceSetup.DeviceID, &deviceSetup.Pins, &deviceSetup.CreatedAt); err != nil {
		return nil, err
	}

	return deviceSetup, nil
}
