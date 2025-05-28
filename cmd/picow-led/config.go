package main

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"picow-led/internal/database"
	"picow-led/internal/micro"
	"sync"

	"gopkg.in/yaml.v3"
)

type ConfigDevice struct {
	Addr string  `yaml:"addr"`
	Name string  `yaml:"name"`
	Pins []uint8 `yaml:"pins"`
}

type Config struct {
	Devices []ConfigDevice `yaml:"devices"`
}

func (c *Config) GetDataBaseDevices() []*database.Device {
	devices := []*database.Device{}

	if c.Devices == nil {
		return devices
	}

	wg := &sync.WaitGroup{}
	for _, d := range c.Devices {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if d.Pins != nil {
				if len(d.Pins) > 0 {
					if err := micro.SetPins(d.Addr, d.Pins); err != nil {
						slog.Warn("Set pins failed", "error", err,
							"device.address", d.Addr, "device.name", d.Name)
					}
				}
			}

			device := database.NewDevice()
			devices = append(devices, device)

			pins, err := micro.GetPins(d.Addr)
			if err != nil {
				slog.Warn("Get Pins failed", "error", err,
					"device.address", d.Addr, "device.name", d.Name)
				device.Error = append(device.Error, err.Error())
				return
			}

			if pins != nil {
				device.Pins = pins
			}

			if color, err := micro.GetColor(d.Addr); err != nil {
				slog.Warn("Get color failed", "error", err,
					"device.address", d.Addr, "device.name", d.Name)
				device.Error = append(device.Error, err.Error())
			} else {
				device.SetColor(color)
			}
		}()
	}
	wg.Wait()

	return devices
}

func loadConfig(paths ...string) (*Config, error) {
	config := &Config{}

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
			return nil, err
		}
		err = yaml.Unmarshal(d, config)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}
