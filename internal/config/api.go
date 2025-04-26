package config

import (
	"io"
	"os"
	"path/filepath"
	"picow-led/internal/api"

	"gopkg.in/yaml.v3"
)

func GetApiOptions(paths ...string) (*api.Options, error) {
	o := &api.Options{
		Servers: []*api.Server{},
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

	return o, nil
}
