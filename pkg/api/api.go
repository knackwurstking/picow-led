package api

import (
	"fmt"
	"sync"

	"github.com/knackwurstking/picow-led-server/pkg/events"
)

const (
	EventNameDevice  = events.Name("device")
	EventNameDevices = events.Name("devices")
	EventNameColor   = events.Name("color")
	EventNameColors  = events.Name("colors")
)

type API struct {
	Colors  APIColors  `json:"colors"`
	Devices APIDevices `json:"devices"`

	EventDeviceChange  *events.Event[*Device]          `json:"-"`
	EventDevicesChange *events.Event[[]*Device]        `json:"-"`
	EventColorChange   *events.Event[ColorEvent]       `json:"-"`
	EventColorsChange  *events.Event[map[string]Color] `json:"-"`

	devicesMutex *sync.Mutex
	colorsMutex  *sync.Mutex
}

func NewAPI() *API {
	return &API{
		Devices: make(APIDevices, 0),
		Colors:  make(APIColors, 0),

		EventDeviceChange:  events.NewEvent[*Device](EventNameDevice),
		EventDevicesChange: events.NewEvent[[]*Device](EventNameDevices),
		EventColorChange:   events.NewEvent[ColorEvent](EventNameColor),
		EventColorsChange:  events.NewEvent[map[string]Color](EventNameColors),

		devicesMutex: &sync.Mutex{},
		colorsMutex:  &sync.Mutex{},
	}
}

func (api *API) HasDevice(addr string) bool {
	for _, d := range api.Devices {
		if d.Server.Addr == addr {
			return true
		}
	}
	return false
}

func (api *API) AddDevice(device *Device) error {
	defer api.devicesMutex.Unlock()
	api.devicesMutex.Lock()

	for _, d := range api.Devices {
		if d.Server.Addr == device.Server.Addr {
			return fmt.Errorf("device \"%s\" already exists", device.Server.Addr)
		}
	}

	defer func() {
		go api.EventDevicesChange.Dispatch(api.Devices)
	}()

	api.Devices = append(api.Devices, device)
	return device.Sync()
}

func (api *API) UpdateDevice(device *Device) error {
	defer api.devicesMutex.Unlock()
	api.devicesMutex.Lock()

	for i, d := range api.Devices {
		if d.Server.Addr == device.Server.Addr {
			defer func() {
				go api.EventDeviceChange.Dispatch(device)
			}()

			api.Devices[i] = device
			return device.Sync()
		}
	}

	return fmt.Errorf("device \"%s\" not found", device.Server.Addr)
}

func (api *API) DeleteDevice(addr string) {
	defer api.devicesMutex.Unlock()
	api.devicesMutex.Lock()

	newDevices := make([]*Device, 0)
	for _, d := range api.Devices {
		if d.Server.Addr == addr {
			defer func() {
				go api.EventDevicesChange.Dispatch(api.Devices)
			}()

			d.Server.Close()
			break
		}

		newDevices = append(newDevices, d)
	}
	api.Devices = newDevices
}

func (api *API) UpdateDevicePins(addr string, pins Pins) error {
	for _, d := range api.Devices {
		if d.Server.Addr == addr {
			defer func() {
				go api.EventDeviceChange.Dispatch(d)
			}()

			return d.SetPins(pins)
		}
	}

	return fmt.Errorf("device %s not found", addr)
}

func (api *API) UpdateDeviceColor(addr string, color Color) error {
	for _, d := range api.Devices {
		if d.Server.Addr == addr {
			defer func() {
				go api.EventDeviceChange.Dispatch(d)
			}()

			return d.SetColor(color)
		}
	}

	return fmt.Errorf("device %s not found", addr)
}

func (api *API) AddColor(name string, color Color) error {
	defer api.colorsMutex.Unlock()
	api.colorsMutex.Lock()

	if _, ok := api.Colors[name]; ok {
		return fmt.Errorf("color \"%s\" already exists", name)
	}

	defer func() {
		go api.EventColorsChange.Dispatch(api.Colors)
	}()

	api.Colors[name] = color

	return nil
}

func (api *API) ReplaceColor(name string, color Color) error {
	defer api.colorsMutex.Unlock()
	api.colorsMutex.Lock()

	if _, ok := api.Colors[name]; !ok {
		return fmt.Errorf("color \"%s\" not exists", name)
	}

	defer func() {
		go api.EventColorChange.Dispatch(NewColorEvent(name, color))
	}()

	api.Colors[name] = color

	return nil
}

func (api *API) DeleteColor(name string) {
	defer api.colorsMutex.Unlock()
	api.colorsMutex.Lock()

	if _, ok := api.Colors[name]; ok {
		defer func() {
			go api.EventColorsChange.Dispatch(api.Colors)
		}()

		delete(api.Colors, name)
	}
}
