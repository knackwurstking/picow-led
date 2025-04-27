package api

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type Device struct {
	Server *Server `json:"server"`

	Online bool   `json:"online" yaml:"-"` //  Not used in configurations
	Error  string `json:"error" yaml:"-"`  //  Not used in configurations

	// Color can be nil
	Color MicroColor `json:"color"`
	// Pins can be nil
	Pins MicroPins `json:"pins"`
}

func GetApiConfig(paths ...string) (*Config, error) {
	o := &Config{
		Servers: []*Server{},
	}

	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			absPath = path
		}
		f, err := os.Open(absPath)
		if err != nil {
			continue
		}
		d, err := io.ReadAll(f)
		if err != nil {
			return o, err
		}
		err = yaml.Unmarshal(d, o)
		if err != nil {
			return o, err
		}
	}

	// TODO: On this point trigger and api update to initially get all data

	return o, nil
}

// TODO: Need to store all data somewhere in the memory
func GetDevices(o *Config) []*Device {
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

			if pins, err := r.Pins(device); err != nil && !device.Online {
				return
			} else {
				device.Pins = pins
			}

			if color, err := r.Color(device); err != nil && !device.Online {
				return
			} else {
				device.Color = color
			}
		}()
	}
	wg.Wait()

	return devices
}
