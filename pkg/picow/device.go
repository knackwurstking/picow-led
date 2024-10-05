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

// NOTE: No error check on "SetPins" and "SetColor" methods, picow server id set to -1
//
// TODO: Missing "GetPins" and "GetColor" methods
type Device struct {
	socket    net.Conn   `json:"-"`
	data      DeviceData `json:"-"`
	connected bool       `json:"-"`
}

func NewDevice(data DeviceData) *Device {
	d := &Device{}
	d.SetData(data)
	return d
}

func (d *Device) Addr() string {
	return d.data.Server.Addr
}

func (d *Device) SetData(data DeviceData) {
	d.data = data

	// Set color and pins to device, ignore any error?
	_ = d.SetPins(d.data.Pins)
	_ = d.SetColor(d.data.Color)
}

func (d *Device) SetPins(p Pins) error {
	if p == nil {
		panic("pins should not be nil")
	}

	if !d.IsConnected() {
		if err := d.Connect(); err != nil {
			return err
		}
		defer d.Close()
	}

	slog.Debug("Set device pins",
		"device.address", d.socket.RemoteAddr(),
		"pins", p,
	)

	req := &Request{
		Type:    "set",
		Group:   "config",
		Command: "led",
		Args:    p.StringArray(),
		ID:      IDNoResponse,
	}

	data, _ := json.Marshal(req)
	_, err := d.socket.Write(data)
	if err == nil {
		d.data.Pins = p
	}
	return err
}

func (d *Device) SetColor(c Color) error {
	if c == nil {
		panic("color should not be nil")
	}

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

	if d.data.Server.Addr == "" {
		return nil
	}

	if d.data.Pins != nil {
		if err := d.SetPins(d.data.Pins); err != nil {
			return err
		}
	}

	if d.data.Color != nil {
		if err := d.SetColor(d.data.Color); err != nil {
			return err
		}
	}

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
