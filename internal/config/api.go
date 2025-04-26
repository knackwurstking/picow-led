package config

import "picow-led/internal/api"

func GetApiOptions(paths ...string) *api.Options {
	o := &api.Options{
		Servers: []*api.Server{},
	}

	// TODO: Read configuration from path if possible

	return o
}
