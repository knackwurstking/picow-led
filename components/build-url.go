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

func HxUrlHomeSectionDevices() templ.SafeURL {
	return templ.SafeURL(BuildUrl("/htmx/home/section/devices"))
}

func HxUrlHomeSectionGroups() templ.SafeURL {
	return templ.SafeURL(BuildUrl("/htmx/home/section/groups"))
}

func HxUrlEditDeviceDialog(deviceID *models.DeviceID) templ.SafeURL {
	if deviceID == nil {
		return templ.SafeURL(BuildUrl("/htmx/dialog/edit-device"))
	}

	return templ.SafeURL(BuildUrl(
		fmt.Sprintf("/htmx/dialog/edit-device?id=%d", *deviceID),
	))
}

func HxUrlDeleteDevice(deviceID models.DeviceID) templ.SafeURL {
	return templ.SafeURL(BuildUrl(
		fmt.Sprintf("/htmx/devices/delete?id=%d", deviceID),
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

// TODO: Power toggle
