package config

import "picow-led/internal/api"

func GetApiOptions(paths ...string) *api.Options {
	// TODO: Get this api options from "~/.config/picow-led/api.yml" or
	// 		 "api.yml" (use yaml v3)
	o := &api.Options{}

	o.Servers = []*api.Server{
		// TODO: Add some servers here for testing
	}

	return o
}
