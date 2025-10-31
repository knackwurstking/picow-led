package services

type DeviceControl struct {
	registry *Registry
}

func NewDeviceControl(registry *Registry) *DeviceControl {
	return &DeviceControl{
		registry: registry,
	}
}

func (p *DeviceControl) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS device_control (
		device_id INTEGER UNIQUE NOT NULL,
		color TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := p.registry.db.Exec(query)
	return err
}

// TODO: Methods for read and set current color, always store color not 0 in table. Get current color from the picow device directly

// TODO: Also handle version, temp and disk-usage
