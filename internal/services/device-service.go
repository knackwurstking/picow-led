package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/pkg/models"
	"github.com/knackwurstking/ui"
)

type DeviceService struct {
	registry *Registry

	pinCache sync.Map
	log      *ui.Logger
}

func NewDeviceService(r *Registry) *DeviceService {
	return &DeviceService{
		registry: r,
		log:      env.NewLogger("DeviceService"),
	}
}

func (d *DeviceService) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS devices (
		id INTEGER NOT NULL,
		addr TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		duty TEXT,

		PRIMARY KEY ("id" AUTOINCREMENT)
	);`

	if _, err := d.registry.db.Exec(query); err != nil {
		return fmt.Errorf("create devices table: %w", err)
	}

	return nil
}

func (d *DeviceService) Get(deviceID models.ID) (*models.Device, error) {
	query := `SELECT * FROM devices WHERE id = ?`
	device, err := ScanDevice(d.registry.db.QueryRow(query, deviceID))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorNotFound
		}
		return nil, fmt.Errorf("get device by ID: %w", err)
	}

	return device, nil
}

func (d *DeviceService) GetByAddr(addr string) (*models.Device, error) {
	query := `SELECT * FROM devices WHERE addr = ?`
	device, err := ScanDevice(d.registry.db.QueryRow(query, addr))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorNotFound
		}
		return nil, fmt.Errorf("get device by address: %w", err)
	}

	return device, nil
}

func (d *DeviceService) List() ([]*models.Device, error) {
	rows, err := d.registry.db.Query(`SELECT * FROM devices`)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.Device{}, nil
		}
		return nil, fmt.Errorf("list devices: %w", err)
	}
	defer rows.Close()

	var devices []*models.Device
	for rows.Next() {
		device, err := ScanDevice(rows)
		if err != nil {
			return nil, fmt.Errorf("scan device from rows: %w", err)
		}
		devices = append(devices, device)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate device rows: %w", err)
	}

	return devices, nil
}

func (d *DeviceService) Add(device *models.Device) (models.ID, error) {
	if err := device.Validate(); err != nil {
		return 0, ErrorValidation
	}

	dutyString, _ := json.Marshal(device.Duty)
	query := `INSERT INTO devices (addr, name, type, duty) VALUES (?, ?, ?, ?)`
	result, err := d.registry.db.Exec(query, device.Addr, device.Name, device.Type, dutyString)
	if err != nil {
		return 0, fmt.Errorf("add device: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last inserted device ID: %w", err)
	}

	return models.ID(id), nil
}

func (d *DeviceService) Update(device *models.Device) error {
	if err := device.Validate(); err != nil {
		return ErrorValidation
	}

	dutyString, _ := json.Marshal(device.Duty)
	query := `UPDATE devices SET addr = ?, name = ?, type = ?, duty = ? WHERE id = ?`
	if r, err := d.registry.db.Exec(query, device.Addr, device.Name, device.Type, dutyString, device.ID); err != nil {
		if i, _ := r.RowsAffected(); i == 0 {
			return ErrorNotFound
		}
		return fmt.Errorf("update device: %w", err)
	}

	return nil
}

func (d *DeviceService) Delete(deviceID models.ID) error {
	query := `DELETE FROM devices WHERE id = ?`
	if r, err := d.registry.db.Exec(query, deviceID); err != nil {
		if i, _ := r.RowsAffected(); i == 0 {
			return nil
		}
		return fmt.Errorf("delete device: %w", err)
	}

	return nil
}

func ScanDevice(r Scannable) (*models.Device, error) {
	device := &models.Device{}
	var dutyString string
	err := r.Scan(&device.ID, &device.Addr, &device.Name, &device.Type, &dutyString)
	json.Unmarshal([]byte(dutyString), &device.Duty)
	return device, err
}

var _ Service = (*DeviceService)(nil)
