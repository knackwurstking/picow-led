package models

import (
	"strconv"
	"strings"
	"time"
)

type Device struct {
	ID        DeviceID  `json:"id"`
	Addr      Addr      `json:"addr"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func NewDevice(addr Addr, name string) *Device {
	// Validate address for host:port
	if !strings.Contains(string(addr), ":") {
		panic("invalid address: " + addr)
	}
	s := strings.Split(string(addr), ":")
	if len(s) != 2 {
		panic("invalid address: " + addr)
	}
	_, err := strconv.Atoi(s[1])
	if err != nil {
		panic("invalid port: " + s[1])
	}

	return &Device{
		Addr:      addr,
		Name:      name,
		CreatedAt: time.Now(),
	}
}

func (d *Device) Validate() bool {
	return d.ID >= 0 && d.Addr != ""
}

func (d *Device) GetAddrSplit() (string, int) {
	s := strings.Split(string(d.Addr), ":")
	address := s[0]
	port, _ := strconv.Atoi(s[1])
	return address, port
}

var _ ServiceModel = (*Device)(nil)
