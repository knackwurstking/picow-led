package picow

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

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

	data, err := json.MarshalIndent(a, "", "\t")
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
