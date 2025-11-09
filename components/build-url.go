package components

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/knackwurstking/picow-led/env"
	"github.com/knackwurstking/picow-led/models"
)

func BuildUrl(path string) string {
	return env.Args.ServerPathPrefix + path
}

// **** //
// Home //
// **** //

func HxUrlHomeSectionDevices() templ.SafeURL {
	return templ.SafeURL(BuildUrl("/htmx/home/devices"))
}

func HxUrlDeleteDevice(deviceID models.DeviceID) templ.SafeURL {
	return templ.SafeURL(BuildUrl(
		fmt.Sprintf("/htmx/home/devices/delete?id=%d", deviceID),
	))
}

func HxUrlTogglePowerDevice(deviceID models.DeviceID) templ.SafeURL {
	return templ.SafeURL(BuildUrl(
		fmt.Sprintf("/htmx/home/devices/toggle-power?id=%d", deviceID),
	))
}

func HxUrlHomeSectionGroups() templ.SafeURL {
	return templ.SafeURL(BuildUrl("/htmx/home/groups"))
}

func HxUrlDeleteGroup(groupID models.GroupID) templ.SafeURL {
	return templ.SafeURL(BuildUrl(
		fmt.Sprintf("/htmx/home/groups/delete?id=%d", groupID),
	))
}

func HxUrlTogglePowerGroup(deviceID models.DeviceID) templ.SafeURL {
	return templ.SafeURL(BuildUrl(
		fmt.Sprintf("/htmx/home/groups/toggle-power?id=%d", deviceID),
	))
}

// ******* //
// Dialogs //
// ******* //

func HxUrlEditDeviceDialog(deviceID *models.DeviceID) templ.SafeURL {
	if deviceID == nil {
		return templ.SafeURL(BuildUrl("/htmx/dialog/edit-device"))
	}

	return templ.SafeURL(BuildUrl(
		fmt.Sprintf("/htmx/dialog/edit-device?id=%d", *deviceID),
	))
}

func HxUrlEditGroupDialog(groupID *models.GroupID) templ.SafeURL {
	if groupID == nil {
		return templ.SafeURL(BuildUrl("/htmx/dialog/edit-group"))
	}

	return templ.SafeURL(BuildUrl(
		fmt.Sprintf("/htmx/dialog/edit-group?id=%d", *groupID),
	))
}
