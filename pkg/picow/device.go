package picow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"sync"
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
type Device struct {
	socket    net.Conn    `json:"-"`
	mutex     *sync.Mutex `json:"-"`
	data      DeviceData  `json:"-"`
	connected bool        `json:"-"`
}

func NewDevice(data DeviceData) *Device {
	d := &Device{
		mutex: &sync.Mutex{},
	}

	d.SetData(data)
	return d
}

func (d *Device) Addr() string {
	return d.data.Server.Addr
}

func (d *Device) SetData(data DeviceData) {
	d.data = data

	// Set color and pins to device, ignore any error?
	if d.data.Pins != nil {
		_ = d.SetPins(d.data.Pins)
	}

	if d.data.Color != nil {
		_ = d.SetColor(d.data.Color)
	}
}

func (d *Device) GetPins() (Pins, error) {
	if !d.IsConnected() {
		if err := d.Connect(); err != nil {
			return nil, err
		}
		defer d.Close()
	}

	slog.Debug("Get device pins", "device.address", d.socket.RemoteAddr())

	req := &Request{
		Type:    "get",
		Group:   "config",
		Command: "led",
		Args:    nil,
		ID:      ID(0),
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()
	data, _ := json.Marshal(req)
	_, err := d.socket.Write(data)
	if err != nil {
		return nil, err
	}

	data, err = d.readAllDataWithTimeout()
	if err != nil {
		return nil, err
	}

	// NOTE: Ignore the `ID`
	resp := &Response[Pins]{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	} else if resp.Error != "" {
		return nil, fmt.Errorf("server returned: %s", resp.Error)
	}

	return resp.Data, nil
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

	d.mutex.Lock()
	defer d.mutex.Unlock()
	data, _ := json.Marshal(req)
	_, err := d.socket.Write(data)
	if err == nil {
		d.data.Pins = p
	}
	return err
}

func (d *Device) GetColor() (Color, error) {
	if !d.IsConnected() {
		if err := d.Connect(); err != nil {
			return nil, err
		}
		defer d.Close()
	}

	slog.Debug("Get device color", "device.address", d.socket.RemoteAddr())

	req := &Request{
		Type:    "get",
		Group:   "led",
		Command: "color",
		Args:    nil,
		ID:      ID(0),
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()
	data, _ := json.Marshal(req)
	_, err := d.socket.Write(data)
	if err != nil {
		return nil, err
	}

	data, err = d.readAllDataWithTimeout()
	if err != nil {
		return nil, err
	}

	// NOTE: Ignore the `ID`
	resp := &Response[Color]{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	} else if resp.Error != "" {
		return nil, fmt.Errorf("server returned: %s", resp.Error)
	}

	return resp.Data, nil
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

	d.mutex.Lock()
	defer d.mutex.Unlock()
	data, _ := json.Marshal(req)
	_, err := d.socket.Write(data)
	if err == nil {
		d.data.Color = c
	}
	return err
}

func (d *Device) MarshalJSON() ([]byte, error) {
	// TODO: Get "pins" and "color" from server

	return json.MarshalIndent(d.data, "", "\t")
}

func (d *Device) UnmarshalJSON(data []byte) error {
	if d.mutex == nil {
		d.mutex = &sync.Mutex{}
	}

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

func (d *Device) readAllDataWithTimeout() ([]byte, error) {
	data := make([]byte, 0)

	b := make([]byte, 1)
	var err error
	var n int

	for {
		if n, err = d.socket.Read(b); err != nil {
			return nil, err
		} else if n == 0 {
			return nil, fmt.Errorf("no data from %s", d.Addr())
		} else {
			if bytes.Equal(b, EndByte) {
				break
			}

			data = append(data, b...)
		}
	}

	return data, nil
}
