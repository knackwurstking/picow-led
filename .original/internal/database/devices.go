package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	ActiveColor []int      `json:"active_color"`
	Color       []int      `json:"color"`
	Pins        []int      `json:"pins"`
	Power       PowerState `json:"power"`

	// Not stored inside the database

	Error []string `json:"error"`
}

func NewDevice() *Device {
	return &Device{
		ActiveColor: make([]int, 0),
		Color:       make([]int, 0),
		Pins:        make([]int, 0),
		Error:       make([]string, 0),
	}
}

func (d *Device) GetColorForPowerState(state PowerState) []int {
	color := []int{}

	switch state {
	case PowerStateOFF:
		for range d.Pins {
			color = append(color, 0)
		}

	case PowerStateON:
		if len(d.ActiveColor) > 0 {
			color = d.ActiveColor
		} else {
			for range d.Pins {
				color = append(color, 255)
			}
		}
	}

	return color
}

func (d *Device) SetColor(color []int) {
	if color == nil {
		return
	}

	d.Color = color
	if slices.Max(d.Color) > 0 {
		d.Power = 1
		d.ActiveColor = color
	} else {
		d.Power = 0
	}

	// Handle active color
	if len(d.ActiveColor) == 0 {
		d.ActiveColor = []int{}
		for range d.Pins {
			d.ActiveColor = append(d.ActiveColor, 255)
		}
	}
}

type Devices struct {
	db *sql.DB
}

func NewDevices(db *sql.DB) (*Devices, error) {
	tables := []string{}

	err := func() error {
		// Query table names
		r, err := db.Query(`SELECT name FROM sqlite_master WHERE type = "table" AND name NOT LIKE 'sqlite_%'`)
		if err != nil {
			return err
		}
		defer r.Close()

		// Scan table names
		var name string
		for r.Next() {
			err = r.Scan(&name)
			if err != nil {
				return err
			}
			tables = append(tables, name)
		}

		return nil
	}()
	if err != nil {
		return nil, err
	}

	if !slices.Contains(tables, "devices") {
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
	query := fmt.Sprintf(`SELECT %s FROM devices ORDER BY name;`, deviceQueryKeys)
	r, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer r.Close()

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
		"SELECT %s FROM devices WHERE addr=\"%s\";",
		deviceQueryKeys, addr,
	)
	r, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if !r.Next() {
		return nil, fmt.Errorf("not found")
	}

	return d.scanDevice(r)
}

func (d *Devices) Set(devices ...*Device) error {
	err := d.DeleteAll()
	if err != nil {
		return err
	}

	var query string
	for _, device := range devices {
		query = fmt.Sprintf(
			"INSERT INTO devices (%s) VALUES (\"%s\", \"%s\", :active_color, :color, :pins, %d);\n",
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
		query = fmt.Sprintf(
			"INSERT INTO devices (%s) VALUES (\"%s\", \"%s\", :active_color, :color, :pins, %d);\n",
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
		"INSERT OR REPLACE INTO devices (%s) VALUES (\"%s\", \"%s\", :active_color, :color, :pins, %d);\n",
		deviceQueryKeys, device.Addr, device.Name, device.Power,
	)

	return d.execDevice(query, device)
}

func (c *Devices) DeleteAll() error {
	_, err := c.db.Exec(`DELETE FROM devices`)
	return err
}

func (c *Devices) Delete(addr string) error {
	query := fmt.Sprintf("DELETE FROM devices WHERE addr=\"%s\";", addr)
	_, err := c.db.Exec(query)
	return err
}

func (d *Devices) Close() {
	d.db.Close()
}

func (d *Devices) scanDevice(r *sql.Rows) (*Device, error) {
	device := NewDevice()

	var (
		activeColorJSON []byte
		colorJSON       []byte
		pinsJSON        []byte
	)

	err := r.Scan(&device.Addr, &device.Name,
		&activeColorJSON, &colorJSON, &pinsJSON,
		&device.Power)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(activeColorJSON, &device.ActiveColor)
	_ = json.Unmarshal(colorJSON, &device.Color)
	_ = json.Unmarshal(pinsJSON, &device.Pins)

	return device, err
}

func (d *Devices) execDevice(query string, device *Device) error {
	var (
		activeColorJSON []byte
		colorJSON       []byte
		pinsJSON        []byte
	)

	activeColorJSON, _ = json.Marshal(device.ActiveColor)
	colorJSON, _ = json.Marshal(device.Color)
	pinsJSON, _ = json.Marshal(device.Pins)

	_, err := d.db.Exec(query,
		sql.Named("active_color", activeColorJSON),
		sql.Named("color", colorJSON),
		sql.Named("pins", pinsJSON))
	return err
}
