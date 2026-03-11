package urlb

import "github.com/knackwurstking/picow-led/internal/env"

func Devices() string {
	return env.Route("/htmx/devices")
}

func AddDeviceDialog() string {
	return env.Route("/htmx/dialogs/add-device")
}
