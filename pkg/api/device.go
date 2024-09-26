package api

import (
	"fmt"
	"strconv"

	"github.com/knackwurstking/picow-led-server/pkg/picow"
)

type Pins []uint

func NewPins() Pins {
	return make(Pins, 0)
}

type Device struct {
	Server *picow.Server `json:"server"`
	Pins   Pins          `json:"pins"`
	Color  Color         `json:"color"`
}

func NewDevice(host string, port int) *Device {
	if port == 0 {
		port = picow.Port
	}

	return &Device{
		Server: picow.NewServer(
			fmt.Sprintf("%s:%d", host, port),
		),
	}
}

func (d *Device) Sync() error {
	err := d.SyncUp()
	if err != nil {
		return fmt.Errorf("device sync (up) %s: %s", d.Server.Addr, err)
	}

	err = d.SyncDown()
	if err != nil {
		return fmt.Errorf("device sync (down) %s: %s", d.Server.Addr, err)
	}

	return nil
}

func (d *Device) SyncUp() error {
	if !d.Server.IsConnected {
		err := d.Server.Connect()
		if err != nil {
			return err
		}
		defer d.Server.Close()
	}

	if err := d.SetPins(d.Pins); err != nil {
		return fmt.Errorf("set pins: %s", err)
	}

	if err := d.SetColor(d.Color); err != nil {
		return fmt.Errorf("set color: %s", err)
	}

	return nil
}

func (d *Device) SyncDown() error {
	if !d.Server.IsConnected {
		err := d.Server.Connect()
		if err != nil {
			return err
		}
		defer d.Server.Close()
	}

	if err := d.GetPins(); err != nil {
		return fmt.Errorf("get pins: %s", err)
	}

	if err := d.GetColor(); err != nil {
		return fmt.Errorf("get color: %s", err)
	}

	return nil
}

func (d *Device) GetPins() error {
	if err := d.checkConnection(); err != nil {
		return err
	}
	defer d.Server.Close()

	request := picow.NewRequest(
		picow.ID(0),
		picow.TypeGet,
		picow.GroupConfig,
		"led",
	)
	if err := d.Server.Send(request); err != nil {
		return err
	}

	if r, err := d.Server.GetResponse(); err != nil {
		return err
	} else {
		if r.Error != "" {
			return fmt.Errorf("%s: %s", d.Server.Addr, r.Error)
		}

		switch v := r.Data.(type) {
		case []any:
			d.Pins = make(Pins, 0)
			for _, v := range v {
				n, ok := v.(float64)
				if !ok {
					return fmt.Errorf("unexpected data type from %s: %T", d.Server.Addr, v)
				}
				d.Pins = append(d.Pins, uint(n))
			}
		default:
			return fmt.Errorf("unexpected data type from %s: %T", d.Server.Addr, v)
		}
	}

	return nil
}

func (d *Device) SetPins(p Pins) error {
	if p == nil {
		d.Pins = make(Pins, 0)
	} else {
		d.Pins = p
	}

	if err := d.checkConnection(); err != nil {
		return err
	}
	defer d.Server.Close()

	pins := make([]string, 0)
	for _, n := range p {
		pins = append(pins, strconv.Itoa(int(n)))
	}

	err := d.Server.Send(
		picow.NewRequest(
			picow.IDNoResponse,
			picow.TypeSet,
			picow.GroupConfig,
			"led",
			pins...,
		),
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *Device) GetColor() error {
	if err := d.checkConnection(); err != nil {
		return err
	}
	defer d.Server.Close()

	request := picow.NewRequest(
		picow.ID(0),
		picow.TypeGet,
		picow.GroupLED,
		"color",
	)
	if err := d.Server.Send(request); err != nil {
		return err
	}

	if r, err := d.Server.GetResponse(); err != nil {
		return err
	} else {
		if r.Error != "" {
			return fmt.Errorf("%s: %s", d.Server.Addr, r.Error)
		}

		switch v := r.Data.(type) {
		case []any:
			d.Color = make(Color, 0)
			for _, v := range v {
				n, ok := v.(float64)
				if !ok {
					return fmt.Errorf("unexpected data type from %s: %T", d.Server.Addr, v)
				}
				d.Color = append(d.Color, uint(n))
			}
		default:
			return fmt.Errorf("unexpected data type from %s: %T", d.Server.Addr, v)
		}
	}

	return nil
}

func (d *Device) SetColor(c Color) error {
	if c == nil {
		d.Color = make(Color, 0)
	} else {
		d.Color = c
	}

	if err := d.checkConnection(); err != nil {
		return err
	}
	defer d.Server.Close()

	color := make([]string, 0)
	for _, n := range c {
		color = append(color, strconv.Itoa(int(n)))
	}

	err := d.Server.Send(
		picow.NewRequest(
			picow.IDNoResponse,
			picow.TypeSet,
			picow.GroupLED,
			"color",
			color...,
		),
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *Device) checkConnection() error {
	if !d.Server.IsConnected {
		err := d.Server.Connect()
		if err != nil {
			return err
		}
	}

	return nil
}
