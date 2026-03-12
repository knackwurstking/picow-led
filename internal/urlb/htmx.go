package urlb

import (
	"fmt"

	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/pkg/models"
)

func Devices() string {
	return env.Route("/htmx/devices")
}

func AddDeviceDialog() string {
	return env.Route("/htmx/dialogs/add-device")
}

func EditDeviceDialog(id models.ID) string {
	return env.Route(fmt.Sprintf("/htmx/dialogs/edit-device?id=%d", id))
}
