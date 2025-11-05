package services

import (
	"log/slog"
	"sync"

	"github.com/knackwurstking/picow-led/models"
)

func ResolveDevices(r *Registry, devices ...*models.Device) ([]*models.ResolvedDevice, error) {
	resolvedDevices := make([]*models.ResolvedDevice, 0, len(devices))

	wg := &sync.WaitGroup{}
	for _, device := range devices {
		wg.Go(func() {
			color, err := r.DeviceControls.GetCurrentColor(device.ID)
			if err != nil {
				slog.Error("Failed to get current color for device", "id", device.ID, "error", err)
			}
			resolvedDevices = append(resolvedDevices, models.NewResolvedDevice(device, color))
		})
	}
	wg.Wait()

	return resolvedDevices, nil
}

func ResolveGroups(r *Registry, groups ...*models.Group) ([]*models.ResolvedGroup, error) {
	resolvedGroups := make([]*models.ResolvedGroup, 0, len(groups))

	for _, group := range groups {
		var devices []*models.Device

		for _, deviceID := range group.Devices {
			device, err := r.Devices.Get(deviceID)
			if err != nil && !IsNotFoundError(err) {
				return nil, err
			}

			devices = append(devices, device)
		}

		resolvedGroups = append(resolvedGroups, models.NewResolvedGroup(group, devices...))
	}

	return resolvedGroups, nil
}
