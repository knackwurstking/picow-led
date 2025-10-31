package models

type ResolvedDevice struct {
	*Device
	Setup *DeviceSetup

	// TODO: Do i need to add the DeviceControl type here?
}

func NewResolvedDevice(device *Device, setup *DeviceSetup) *ResolvedDevice {
	return &ResolvedDevice{
		Device: device,
		Setup:  setup,
	}
}

type ResolvedGroup struct {
	*Group
	Devices []*ResolvedDevice
}

func NewResolvedGroup(group *Group, devices ...*ResolvedDevice) *ResolvedGroup {
	return &ResolvedGroup{
		Group:   group,
		Devices: devices,
	}
}

var _ ServiceModel = (*ResolvedDevice)(nil)
var _ ServiceModel = (*ResolvedGroup)(nil)
