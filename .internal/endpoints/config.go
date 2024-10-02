package endpoints

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/knackwurstking/picow-led-server/pkg/api"
)

type Config struct {
	Path string
	API  *api.API

	mutex *sync.Mutex
}

func NewConfig(a *api.API) *Config {
	path, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Sprintf("Get user config directory: %s", err))
	}

	path = filepath.Join(path, "picow-led-server", "api.json")

	return &Config{
		Path:  path,
		API:   a,
		mutex: &sync.Mutex{},
	}
}

func (c *Config) save() error {
	data, err := json.Marshal(c.API)
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
		err = json.Unmarshal(data, c.API)
		if err != nil {
			c.mutex.Unlock()
			return err
		}
	}
	c.mutex.Unlock()

	wg := &sync.WaitGroup{}
	for _, d := range c.API.Devices {
		wg.Add(1)
		go func(d *api.Device) {
			defer wg.Done()
			if err := d.SyncUp(); err != nil {
				slog.Error("Sync device", "error", err.Error())
			}
		}(d)
	}
	wg.Wait()

	return nil
}
