package api

import (
	"fmt"
	"sync"
)

var (
	apiDevicesMutex = &sync.Mutex{}
)

type APIDevices []*Device

func NewAPIDevices() APIDevices {
	return make(APIDevices, 0)
}

func (ad *APIDevices) Add(device *Device) error {
	defer apiDevicesMutex.Unlock()
	apiDevicesMutex.Lock()

	for _, d := range *ad {
		if d.Server.Addr == device.Server.Addr {
			return fmt.Errorf("device \"%s\" already exists", device.Server.Addr)
		}
	}

	*ad = append(*ad, device)
	return device.Sync()
}

func (ad APIDevices) Update(device *Device) error {
	defer apiDevicesMutex.Unlock()
	apiDevicesMutex.Lock()

	for i, d := range ad {
		if d.Server.Addr == device.Server.Addr {
			ad[i] = device
			return device.Sync()
		}
	}

	return fmt.Errorf("device \"%s\" not found", device.Server.Addr)
}

func (ad *APIDevices) Remove(addr string) {
	defer apiDevicesMutex.Unlock()
	apiDevicesMutex.Lock()

	nd := make([]*Device, 0)
	for _, d := range *ad {
		if d.Server.Addr != addr {
			nd = append(nd, d)
		} else {
			d.Server.Close()
		}
	}
	*ad = nd
}
