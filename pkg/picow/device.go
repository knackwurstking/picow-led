package picow

import (
	"encoding/json"
	"net"
	"time"
)

type DeviceData_Server struct {
	Name   string `json:"name"`
	Addr   string `json:"addr"`
	Online bool   `json:"online"`
}

type DeviceData struct {
	Server DeviceData_Server `json:"server"`
	Pins   []uint            `json:"pins"`
	Color  []uint            `json:"color"`
}

type Device struct {
	socket    net.Conn   `json:"-"`
	data      DeviceData `json:"-"`
	connected bool       `json:"-"`
}

func (d *Device) Addr() string {
	return d.data.Server.Addr
}

func (d *Device) SetColor(c []uint) error {
	if !d.IsConnected() {
		if err := d.Connect(); err != nil {
			return err
		}
		defer d.Close()
	}

	// TODO: ...

	return nil
}

func (d *Device) MarshalJSON() ([]byte, error) {
	// TODO: Get "pins" and "color" from server

	return json.Marshal(d.data)
}

func (d *Device) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &d.data)
	if err != nil {
		return err
	}

	// TODO: Sync device "pins" and "color"

	return nil
}

func (d *Device) Connect() error {
	dialer := net.Dialer{
		Timeout: time.Duration(time.Second * 5),
	}
	conn, err := dialer.Dial("tcp", d.data.Server.Addr)
	if err != nil {
		d.data.Server.Online = false
		d.connected = false
		return err
	}

	d.connected = true
	d.data.Server.Online = true
	d.socket = conn

	return nil
}

func (d *Device) Send() error {
	if !d.IsConnected() {
		if err := d.Connect(); err != nil {
			return err
		}
		defer d.Close()
	}

	// TODO: ...

	return nil
}

func (d *Device) Read() error {
	if !d.IsConnected() {
		if err := d.Connect(); err != nil {
			return err
		}
		defer d.Close()
	}

	// TODO: ...

	return nil
}

func (d *Device) IsConnected() bool {
	return d.connected
}

func (d *Device) IsOnline() bool {
	return d.data.Server.Online
}

func (d *Device) Close() {
	d.socket.Close()
	d.connected = false
}
