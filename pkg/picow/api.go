package picow

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

var devicesMutex = &sync.Mutex{}

type Devices []*Device

func (d *Devices) Delete(device *Device) {
	devicesMutex.Lock()
	defer devicesMutex.Unlock()

	newDevices := make([]*Device, 0)
	for _, d := range *d {
		if d.Addr() != device.Addr() {
			newDevices = append(newDevices, d)
		}
	}
	*d = newDevices
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

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, a)
}

func (a *Api) SaveToPath(path string) error {
	// TODO: Create dirs to path, strip last index "/" (file name)

	// TODO: Marshal data and write to file

	return nil
}
