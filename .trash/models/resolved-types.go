package models

type ResolvedDevice struct {
	*Device
	CurrentColor []uint8 `json:"color"`
}

func NewResolvedDevice(device *Device, currentColor []uint8) *ResolvedDevice {
	return &ResolvedDevice{
		Device:       device,
		CurrentColor: currentColor,
	}
}

type ResolvedGroup struct {
	*Group
	Devices []*Device `json:"devices"`
}

func NewResolvedGroup(group *Group, devices ...*Device) *ResolvedGroup {
	return &ResolvedGroup{
		Group:   group,
		Devices: devices,
	}
}

var _ ServiceModel = (*ResolvedGroup)(nil)
