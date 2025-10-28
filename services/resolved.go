package services

import "github.com/knackwurstking/picow-led/models"

func ResolveDevices(r *Registry, devices ...*models.Device) ([]*models.ResolvedDevice, error) {
	resolvedDevices := make([]*models.ResolvedDevice, 0, len(devices))
	for _, device := range devices {
		setup, err := r.DeviceSetups.Get(device.ID)
		if err != nil {
			return nil, err
		}
		resolvedDevice := models.NewResolvedDevice(device, setup)
		resolvedDevices = append(resolvedDevices, resolvedDevice)
	}
	return resolvedDevices, nil
}

func ResolveGroups(r *Registry, groups ...*models.Group) ([]*models.ResolvedGroup, error) {
	resolvedGroups := make([]*models.ResolvedGroup, 0, len(groups))

	for _, group := range groups {
		var devices []*models.ResolvedDevice

		for _, deviceID := range group.Devices {
			device, err := r.Devices.Get(deviceID)
			if err != nil {
				return nil, err
			}

			resolvedDevice, err := ResolveDevices(r, device)
			if err != nil {
				return nil, err
			}

			devices = append(devices, resolvedDevice...)
		}

		resolvedGroups = append(resolvedGroups, models.NewResolvedGroup(group, devices...))
	}

	return resolvedGroups, nil
}
