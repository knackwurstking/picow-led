package urlb

import (
	"fmt"
	"net/url"

	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/pkg/models"
)

type PowerMode string

const (
	PowerModeOn  PowerMode = "on"
	PowerModeOff PowerMode = "off"
)

func Devices() string {
	return env.Route("/htmx/devices")
}

func ToggleDevicePower(id models.ID) string {
	u := url.URL{}
	u.Path = env.Route("/htmx/devices/toggle-power")

	q := u.Query()
	q.Set("id", fmt.Sprintf("%d", id))
	u.RawQuery = q.Encode()

	return u.String()
}

func AddDeviceDialog() string {
	return env.Route("/htmx/dialogs/add-device")
}

func EditDeviceDialog(id models.ID) string {
	u := url.URL{}
	u.Path = env.Route("/htmx/dialogs/edit-device")

	q := u.Query()
	q.Set("id", fmt.Sprintf("%d", id))
	u.RawQuery = q.Encode()

	return u.String()
}

func Groups() string {
	return env.Route("/htmx/groups")
}

func PowerGroup(id models.ID, mode PowerMode) string {
	u := url.URL{}
	u.Path = env.Route("/htmx/groups/power")

	q := u.Query()
	q.Set("id", fmt.Sprintf("%d", id))
	q.Set("mode", string(mode))
	u.RawQuery = q.Encode()

	return u.String()
}

func AddGroupDialog() string {
	return env.Route("/htmx/dialogs/add-group")
}

func EditGroupDialog(groupID models.ID) string {
	u := url.URL{}
	u.Path = env.Route("/htmx/dialogs/edit-group")

	q := u.Query()
	q.Set("id", fmt.Sprintf("%d", groupID))
	u.RawQuery = q.Encode()

	return u.String()
}
