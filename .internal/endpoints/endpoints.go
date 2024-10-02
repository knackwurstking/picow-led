package endpoints

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"

	"github.com/knackwurstking/picow-led-server/pkg/api"
	"github.com/knackwurstking/picow-led-server/pkg/clients"
)

const (
	ClientsEmitTypeDevice  = clients.EmitType("device")
	ClientsEmitTypeDevices = clients.EmitType("devices")
	ClientsEmitTypeColor   = clients.EmitType("color")
	ClientsEmitTypeColors  = clients.EmitType("colors")
)

func Create(e *echo.Echo) {
	c := clients.NewClients()
	a := api.NewAPI()

	_ = a.EventDeviceChange.Add("emitter", func(d *api.Device) {
		c.Emit(ClientsEmitTypeDevice, d)
	})

	_ = a.EventDevicesChange.Add("emitter", func(d []*api.Device) {
		c.Emit(ClientsEmitTypeDevices, d)
	})

	_ = a.EventColorChange.Add("emitter", func(ce api.ColorEvent) {
		c.Emit(ClientsEmitTypeColor, ce)
	})

	_ = a.EventColorsChange.Add("emitter", func(mc map[string]api.Color) {
		c.Emit(ClientsEmitTypeColors, mc)
	})

	config := NewConfig(a)
	slog.Debug("Try to load API configuration", "path", config.Path)
	if err := config.load(); err != nil {
		panic(fmt.Sprintf("Loading API configuration froom \"%s\": %s", config.Path, err))
	}

	createStaticEndpoints(e)
	createEventsEndpoints(e, c)

	createApiEndpoints(e, a, func() {
		if err := config.save(); err != nil {
			slog.Warn("Save API configuration", "path", config.Path, "error", err)
		}
	})
}
