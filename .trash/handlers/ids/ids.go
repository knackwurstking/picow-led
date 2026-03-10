package ids

import (
	"fmt"

	"github.com/knackwurstking/picow-led/models"
)

const (
	HomeSectionDevices        string = "section-devices"
	HomeSectionGroups         string = "section-groups"
	HomeSectionDevicesList    string = "devices-list"
	HomeSectionDevicesSpinner string = "section-devices-spinner"
	HomeSectionGroupsSpinner  string = "section-groups-spinner"

	DialogEditGroup    string = "edit-group-dialog"
	DialogEditDevice   string = "edit-device-dialog"
	DialogNewDevice    string = "new-device-dialog"
	DialogNewGroup     string = "new-group-dialog"
	DialogGroupDevices string = "group-devices"
)

func HomeSectionDevicesDevice(id models.DeviceID) string {
	return fmt.Sprintf("device-%d", id)
}

func HomeSectionGroupsGroup(id models.GroupID) string {
	return fmt.Sprintf("group-%d", id)
}

func HomeSectionDevicesDeviceError(id models.DeviceID) string {
	return fmt.Sprintf("device-%d-error", id)
}

func HomeSectionGroupsGroupError(id models.GroupID) string {
	return fmt.Sprintf("group-%d-error", id)
}

// Dialogs

func DialogDeviceCheckbox(id models.DeviceID) string {
	return fmt.Sprintf("group-device-checkbox-%d", id)
}

// Home: Power Button

func HomeSectionDevicesPowerButton(id models.DeviceID) string {
	return fmt.Sprintf("device-%d-power-button", id)
}

func HomeSectionGroupsPowerButton(id models.GroupID) string {
	return fmt.Sprintf("group-%d-power-button", id)
}
