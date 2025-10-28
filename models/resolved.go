package models

type ResolvedDevice struct {
	*Device
	Setup *DeviceSetup
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
