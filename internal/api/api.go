package api

import "sync"

type Device struct {
	Server *Server `json:"server"`
	// Color can be nil
	Color MicroColor `json:"color"`
	// Pins can be nil
	Pins MicroPins `json:"pins"`
}

func GetDevices(o *Options) []*Device {
	devices := []*Device{}

	for _, server := range o.Servers {
		d := &Device{}
		d.Server = server
		devices = append(devices, d)
	}

	wg := &sync.WaitGroup{}
	for _, device := range devices {
		wg.Add(1)
		go func() {
			defer wg.Done()

			r := &MicroRequest{}

			if pins, err := r.Pins(device.Server); err != nil {
				device.Server.Error = err.Error()
				return
			} else {
				device.Pins = pins
			}

			if color, err := r.Color(device.Server); err != nil {
				device.Server.Error = err.Error()
				return
			} else {
				device.Color = color
			}
		}()
	}
	wg.Wait()

	return devices
}
