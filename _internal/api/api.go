package api

import (
	"log/slog"
	"picow-led/internal/types"
	"sync"
)

func GetDevices(devices ...*types.Device) []*types.Device {
	wg := &sync.WaitGroup{}
	for _, device := range devices {
		wg.Add(1)
		go func() {
			defer wg.Done()

			r := types.NewMicroRequest(types.MicroIDDefault)

			if pins, err := r.Pins(device); err != nil && !device.Online {
				return
			} else {
				device.Pins = pins
			}

			if device.Error != "" {
				slog.Error("Request pins", "error", device.Error, "device.server", device.Server)
			}
			pinsError := device.Error

			if color, err := r.Color(device); err != nil && !device.Online {
				return
			} else {
				device.Color = color
			}

			if device.Error != "" {
				slog.Error("Request color", "error", device.Error, "device.server", device.Server)
			} else {
				device.Error = pinsError
			}
		}()
	}
	wg.Wait()

	return devices
}

func SetColor(c types.MicroColor, devices ...*types.Device) []*types.Device {
	wg := &sync.WaitGroup{}
	for _, d := range devices {
		wg.Add(1)
		go func() {
			defer wg.Done()

			r := types.NewMicroRequest(types.MicroIDDefault)
			if err := r.SetColor(d, c); err != nil {
				d.Error = err.Error()
			} else {
				d.Color = c
			}
		}()
	}
	wg.Wait()

	return devices
}

func SetPins(p types.MicroPins, devices ...*types.Device) []*types.Device {
	wg := &sync.WaitGroup{}
	for _, d := range devices {
		wg.Add(1)
		go func() {
			defer wg.Done()

			r := types.NewMicroRequest(types.MicroIDDefault)
			if err := r.SetPins(d, p); err != nil {
				d.Error = err.Error()
			} else {
				d.Pins = p
			}
		}()
	}
	wg.Wait()

	return devices
}
