package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
)

const (
	PowerStateOFF = 0
	PowerStateON  = 1

	deviceQueryKeys = "addr, name, active_color, color, pins, power"
)

type PowerState uint8

type Device struct {
	Addr        string     `json:"addr"`
	Name        string     `json:"name"`
	ActiveColor []uint8    `json:"active_color"`
	Color       []uint8    `json:"color"`
	Pins        []uint8    `json:"pins"`
	Power       PowerState `json:"power"`

	// Not stored inside the database

	Online bool   `json:"online"`
	Error  string `json:"error"`
}

func NewDevice() *Device {
	return &Device{
		ActiveColor: make([]uint8, 0),
		Color:       make([]uint8, 0),
		Pins:        make([]uint8, 0),
	}
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
				active_color BLOB NOT NULL,
				color BLOB NOT NULL,
				pins BLOB NOT NULL,
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
	query := fmt.Sprintf(`SELECT %s FROM devices`, deviceQueryKeys)
	r, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}

	var (
		devices = []*Device{}
		device  *Device
	)
	for r.Next() {
		device, err = d.scanDevice(r)
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
		deviceQueryKeys, addr,
	)
	r, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}

	r.Next()
	return d.scanDevice(r)
}

func (d *Devices) Set(devices ...*Device) error {
	err := d.DeleteAll()
	if err != nil {
		return err
	}

	var query string
	for _, device := range devices {
		query += fmt.Sprintf(
			"INSERT INTO devices (%s) VALUES (%s, %s, ?, ?, ?, %d);\n",
			deviceQueryKeys, device.Addr, device.Name, device.Power,
		)

		err = d.execDevice(query, device)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Devices) Add(devices ...*Device) error {
	var (
		query string
		err   error
	)
	for _, device := range devices {
		query += fmt.Sprintf(
			"INSERT INTO devices (%s) VALUES (%s, %s, ?, ?, ?, %d);\n",
			deviceQueryKeys, device.Addr, device.Name, device.Power,
		)

		err = d.execDevice(query, device)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Devices) Update(addr string, device *Device) error {
	query := fmt.Sprintf(
		"INSERT OR REPLACE INTO devices (%s) VALUES (%s %s, ?, ?, ?, %d);\n",
		deviceQueryKeys, device.Addr, device.Name, device.Power,
	)

	return d.execDevice(query, device)
}

func (c *Devices) DeleteAll() error {
	_, err := c.db.Exec(`DELETE FROM devices`)
	return err
}

func (c *Devices) Delete(addr string) error {
	query := fmt.Sprintf("DELETE FROM devices WHERE addr=%s;", addr)
	_, err := c.db.Exec(query)
	return err
}

func (d *Devices) Close() {
	d.db.Close()
}

func (d *Devices) scanDevice(r *sql.Rows) (*Device, error) {
	device := &Device{}

	var (
		activeColorJSON []byte
		colorJSON       []byte
		pinsJSON        []byte
	)

	err := r.Scan(&device.Addr, &device.Name,
		activeColorJSON, colorJSON, pinsJSON,
		&device.Power)
	if err != nil {
		return nil, err
	}

	device.ActiveColor = activeColorJSON
	device.Color = colorJSON
	device.Pins = pinsJSON

	return device, err
}

func (d *Devices) execDevice(query string, device *Device) error {
	var (
		activeColorJSON []byte
		colorJSON       []byte
		pinsJSON        []byte
	)

	activeColorJSON, _ = json.Marshal(device.Color)
	colorJSON, _ = json.Marshal(device.Color)
	pinsJSON, _ = json.Marshal(device.Color)

	_, err := d.db.Exec(query, activeColorJSON, colorJSON, pinsJSON)
	return err
}
