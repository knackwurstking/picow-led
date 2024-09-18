package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	_api "github.com/knackwurstking/picow-led-server/pkg/api"
)

type Config struct {
	Path string

	mutex *sync.Mutex
}

func NewConfig() *Config {
	path, err := os.UserConfigDir()
	if err != nil {
		slog.Error("Get user config directory", "err", err)
		os.Exit(ErrorCodeConfiguration)
	}

	path = filepath.Join(path, "picow-led-server", "api.json")

	return &Config{
		Path:  path,
		mutex: &sync.Mutex{},
	}
}

func (c *Config) save() error {
	data, err := json.Marshal(api)
	if err != nil {
		return err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	return os.WriteFile(c.Path, data, os.FileMode(0644))
}

func (c *Config) load() error {
	c.mutex.Lock()
	data, err := os.ReadFile(c.Path)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, api)
		if err != nil {
			c.mutex.Unlock()
			return err
		}
	}
	c.mutex.Unlock()

	wg := &sync.WaitGroup{}
	for _, d := range api.Devices {
		wg.Add(1)
		go func(d *_api.Device) {
			defer wg.Done()
			if err := d.SyncUp(); err != nil {
				slog.Error("Sync device", "err", err.Error())
			}
		}(d)
	}
	wg.Wait()

	return nil
}
