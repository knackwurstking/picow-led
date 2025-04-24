package api

import "sync"

type Device struct {
	Server *Server `json:"server"`
	// Color can be nil
	Color MicroColor `json:"color"`
	// Color can be nil
	Pins MicroPins `json:"pins"`
}

func GetDevices(o Options) []*Device {
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

			// TODO: Request pins and color
			r := &MicroRequest{}

			if pins, err := r.RequestPins(device.Server); err != nil {
				// TODO: Logging...
				return
			} else {
				device.Pins = pins
			}

			if color, err := r.RequestColor(device.Server); err != nil {
				// TODO: Logging...
				return
			} else {
				device.Color = color
			}
		}()
	}
	wg.Wait()

	return devices
}
