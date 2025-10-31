package services

import (
	"github.com/knackwurstking/picow-led/models"
)

func ResolveDevice(r *Registry, device *models.Device) (*models.ResolvedDevice, error) {
	setup, err := r.DeviceSetups.Get(device.ID)
	if err != nil && err != ErrNotFound {
		return nil, err
	}

	resolvedDevice := models.NewResolvedDevice(device, setup)
	return resolvedDevice, nil
}

func ResolveDeviceID(r *Registry, deviceID models.DeviceID) (*models.ResolvedDevice, error) {
	device, err := r.Devices.Get(deviceID)
	if err != nil && err != ErrNotFound {
		return nil, err
	}

	return ResolveDevice(r, device)
}

func ResolveDevices(r *Registry, devices ...*models.Device) ([]*models.ResolvedDevice, error) {
	resolvedDevices := make([]*models.ResolvedDevice, 0, len(devices))

	for _, device := range devices {
		resolvedDevice, err := ResolveDevice(r, device)
		if err != nil {
			return nil, err
		}
		resolvedDevices = append(resolvedDevices, resolvedDevice)
	}

	return resolvedDevices, nil
}

func ResolveDevicesID(r *Registry, devicesID ...models.DeviceID) ([]*models.ResolvedDevice, error) {
	resolvedDevices := make([]*models.ResolvedDevice, 0, len(devicesID))

	for _, deviceID := range devicesID {
		resolvedDevice, err := ResolveDeviceID(r, deviceID)
		if err != nil {
			return nil, err
		}
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
			if err != nil && err != ErrNotFound {
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
