package picow

import (
	"slices"
	"sync"
)

type Devices []*Device

func (devices *Devices) Add(device *Device, mutex *sync.Mutex) {
	if device == nil {
		panic("device parameter is nil")
	}

	if mutex != nil {
		mutex.Lock()
		defer mutex.Unlock()
	}

	*devices = append(*devices, NewDevice(*device.DeviceData()))
}

func (devices *Devices) Remove(device *Device, mutex *sync.Mutex) (ok bool) {
	if device == nil {
		panic("device parameter is nil")
	}

	if mutex != nil {
		mutex.Lock()
		defer mutex.Unlock()
	}

	indexToDelete := -1
	for i, d := range *devices {
		if d.Addr() == device.Addr() {
			indexToDelete = i
			break
		}
	}

	if indexToDelete > -1 {
		*devices = slices.Delete[[]*Device](*devices, indexToDelete, indexToDelete)
		ok = true
	}

	return ok
}
