package models

import (
	"errors"
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

func NewDevice(addr Addr, name string) (*Device, error) {
	// Validate address for host:port
	if !strings.Contains(string(addr), ":") {
		return nil, errors.New("invalid address format: " + string(addr) + " (expected host:port)")
	}
	s := strings.Split(string(addr), ":")
	if len(s) != 2 {
		return nil, errors.New("invalid address format: " + string(addr) + " (expected host:port)")
	}
	_, err := strconv.Atoi(s[1])
	if err != nil {
		return nil, errors.New("invalid port in address: " + string(addr) + " (port must be numeric)")
	}

	return &Device{
		Addr:      addr,
		Name:      name,
		CreatedAt: time.Now(),
	}, nil
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
