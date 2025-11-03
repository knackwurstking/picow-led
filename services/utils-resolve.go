package services

import (
	"github.com/knackwurstking/picow-led/models"
)

func ResolveGroups(r *Registry, groups ...*models.Group) ([]*models.ResolvedGroup, error) {
	resolvedGroups := make([]*models.ResolvedGroup, 0, len(groups))

	for _, group := range groups {
		var devices []*models.Device

		for _, deviceID := range group.Devices {
			device, err := r.Devices.Get(deviceID)
			if err != nil && err != ErrNotFound {
				return nil, err
			}

			devices = append(devices, device)
		}

		resolvedGroups = append(resolvedGroups, models.NewResolvedGroup(group, devices...))
	}

	return resolvedGroups, nil
}
