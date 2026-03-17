package urlb

import "github.com/knackwurstking/picow-led/internal/env"

func PageHome() string {
	return env.Route("/")
}

func PageDevice() string {
	return env.Route("/device")
}
