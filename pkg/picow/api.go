package picow

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var devicesMutex = &sync.Mutex{}

type Devices []*Device

func (d *Devices) Delete(device *Device) (ok bool) {
	devicesMutex.Lock()
	defer devicesMutex.Unlock()

	newDevices := make([]*Device, 0)
	for _, d := range *d {
		if d.Addr() != device.Addr() {
			newDevices = append(newDevices, d)
		} else {
			ok = true
		}
	}
	*d = newDevices

	return ok
}

type Api struct {
	Devices Devices `json:"devices"`
}

func NewApi() *Api {
	return &Api{
		Devices: make([]*Device, 0),
	}
}

func (a *Api) LoadFromPath(path string) error {
	i, err := os.Stat(path)
	if err != nil {
		return err
	}

	if i.IsDir() {
		return fmt.Errorf("%s: destination have to be a file", path)
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, a)
}

func (a *Api) SaveToPath(path string) error {
	dir, _ := filepath.Split(path)
	if err := os.MkdirAll(dir, os.FileMode(0700)); err != nil {
		return err
	}

	data, err := json.Marshal(a)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		return err
	}

	n, err := f.Write(data)
	if err != nil {
		return err
	} else if n == 0 {
		return fmt.Errorf("nothing written to %s", path)
	}

	return nil
}
