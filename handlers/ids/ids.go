package ids

import (
	"fmt"

	"github.com/knackwurstking/picow-led/models"
)

const (
	HomeSectionDevices        string = "section-devices"
	HomeSectionDevicesSpinner string = "section-devices-spinner"
	HomeSectionDevicesList    string = "devices-list"
	HomeSectionGroups         string = "section-groups"
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

func HomeSectionDevicesPowerButton(id models.DeviceID) string {
	return fmt.Sprintf("device-%d-power-button", id)
}

func HomeSectionDevicesDeviceError(id models.DeviceID) string {
	return fmt.Sprintf("device-%d-error", id)
}

func HomeSectionGroupsGroup(id models.GroupID) string {
	return fmt.Sprintf("group-%d", id)
}

func HomeSectionGroupsGroupError(id models.GroupID) string {
	return fmt.Sprintf("group-%d-error", id)
}

func DialogDeviceCheckbox(id models.DeviceID) string {
	return fmt.Sprintf("group-device-checkbox-%d", id)
}
