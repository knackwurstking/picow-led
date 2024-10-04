package picow

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Api struct {
	Devices []*Device `json:"devices"`
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
		return fmt.Errorf("%s: destination hav to be a file", path)
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
