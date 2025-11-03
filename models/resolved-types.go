package models

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
