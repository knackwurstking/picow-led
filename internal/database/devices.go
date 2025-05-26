package database

import (
	"database/sql"
	"log/slog"
	"slices"
)

type Device struct {
	// TODO: Continue here
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

	// TODO: create table if not exists
	if !slices.Contains(tables, "devices") {
		slog.Debug("Create (sqlite) database table", "name", "devices")
		_, err = db.Exec(`
      		CREATE TABLE colors (
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

	r, err := d.db.Query(`SELECT ... FROM devices`) // TODO: ...
	if err != nil {
		return nil, err
	}

	var device Device
	for r.Next() {
		err = r.Scan() // TODO: ...
		if err != nil {
			return nil, err
		}
		devices = append(devices, &device)
	}

	return devices, nil
}

func (d *Devices) Close() {
	d.db.Close()
}
