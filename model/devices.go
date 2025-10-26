package model

import (
	"net"
	"time"
)

type DeviceID int64

// Addr contains an IP address and port number separated by a colon.
type Addr string

func (a Addr) SplitHostPort() (host string, port string) {
	host, port, _ = net.SplitHostPort(string(a)) // Ignore error
	return host, port
}

type Device struct {
	ID        DeviceID  `json:"id"`
	Addr      Addr      `json:"addr"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func NewDevice(addr Addr, name string) *Device {
	return &Device{
		Addr:      addr,
		Name:      name,
		CreatedAt: time.Now(),
	}
}

func NewDeviceWithID(id DeviceID, addr Addr, name string, createdAt time.Time) *Device {
	return &Device{
		ID:        id,
		Addr:      addr,
		Name:      name,
		CreatedAt: createdAt,
	}
}

func (d *Device) Validate() bool {
	return d.ID >= 0 && d.Addr != ""
}
