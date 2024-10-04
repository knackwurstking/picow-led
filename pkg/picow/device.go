package picow

import (
	"encoding/json"
	"log/slog"
	"net"
	"time"
)

type (
	DeviceDataServer struct {
		Name   string `json:"name"`
		Addr   string `json:"addr"`
		Online bool   `json:"online"`
	}

	DeviceData struct {
		Server DeviceDataServer `json:"server"`
		Pins   Pins             `json:"pins"`
		Color  Color            `json:"color"`
	}
)

type Device struct {
	socket    net.Conn   `json:"-"`
	data      DeviceData `json:"-"`
	connected bool       `json:"-"`
}

func (d *Device) Addr() string {
	return d.data.Server.Addr
}

func (d *Device) SetColor(c Color) error {
	if !d.IsConnected() {
		if err := d.Connect(); err != nil {
			return err
		}
		defer d.Close()
	}

	slog.Debug("Set device color",
		"device.address", d.socket.RemoteAddr(),
		"color", c,
	)

	req := &Request{
		Type:    "set",
		Group:   "led",
		Command: "color",
		Args:    c.StringArray(),
		ID:      IDNoResponse,
	}

	data, _ := json.Marshal(req)
	_, err := d.socket.Write(data)
	if err == nil {
		d.data.Color = c
	}

	return err
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

func (d *Device) Socket() net.Conn {
	return d.socket
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
