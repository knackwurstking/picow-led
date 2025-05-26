package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"slices"
)

const (
	PowerStateOFF = 0
	PowerStateON  = 1
)

type PowerState uint8

type Device struct {
	Addr        string     `json:"addr"`
	Name        string     `json:"name"`
	Color       []uint8    `json:"color"`
	Pins        []uint8    `json:"pins"`
	ActiveColor []uint8    `json:"active_color"`
	Power       PowerState `json:"power"`
}

type Devices struct {
	db *sql.DB
}

func NewDevices(db *sql.DB) (*Devices, error) {
	// Query table names
	r, err := db.Query(`SELECT name FROM sqlite_master WHERE type = "table" AND name NOT LIKE 'sqlite_%'`)
	if err != nil {
		return nil, err
	}

	// Scan table names
	tables := []string{}
	var name string
	for r.Next() {
		err = r.Scan(&name)
		if err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}

	if !slices.Contains(tables, "devices") {
		slog.Debug("Create (sqlite) database table", "name", "devices")
		_, err = db.Exec(`
      		CREATE TABLE devices (
				addr TEXT NOT NULL,
				name TEXT NOT NULL,
				color BLOB NOT NULL,
				pins BLOB NOT NULL,
				active_color BLOB NOT NULL,
				power INTEGER NOT NULL,
				PRIMARY KEY("addr")
      		);
    	`)
		if err != nil {
			return nil, err
		}
	}

	return &Devices{
		db: db,
	}, nil
}

func (d *Devices) List() ([]*Device, error) {
	devices := []*Device{}

	query := fmt.Sprintf(`SELECT %s FROM devices`, d.deviceQueryKeys())
	r, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}

	var device *Device
	for r.Next() {
		device, err = d.scan(r)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func (d *Devices) Get(addr string) (*Device, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM devices WHERE addr=%s",
		d.deviceQueryKeys(), addr,
	)
	r, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}

	r.Next()
	return d.scan(r)
}

func (d *Devices) Close() {
	d.db.Close()
}

func (d *Devices) scan(r *sql.Rows) (*Device, error) {
	device := &Device{}
	err := r.Scan(&device.Addr, &device.Name,
		&device.Color, &device.Pins, &device.ActiveColor,
		&device.Power)
	if err != nil {
		return nil, err
	}
	return device, err
}

func (d *Devices) deviceQueryKeys() string {
	return "addr, name, color, pins, active_color, power"
}
