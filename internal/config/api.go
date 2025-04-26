package config

import "picow-led/internal/api"

func GetApiOptions(paths ...string) *api.Options {
	o := &api.Options{}

	o.Servers = []*api.Server{
		// TODO: Just testing here, move this to "./api.yaml"
		{Addr: "192.168.178.50:3000", Name: "Living Room"},
		{Addr: "192.168.178.67:3000", Name: "Bed Room"},
	}

	return o
}
