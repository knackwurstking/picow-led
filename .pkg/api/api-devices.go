package api

import (
	"sync"
)

var apiDevicesMutex = &sync.Mutex{}

type APIDevices []*Device

func NewAPIDevices() APIDevices {
	return make(APIDevices, 0)
}
