package routes

import (
	"encoding/json"
	"net/http"
	"picow-led/internal/api"

	"github.com/labstack/echo/v4"
)

type Api struct {
	ServerPathPrefix string
	Config           *api.Config
}

type RequestDevicesColorData struct {
	Devices []*api.Device  `json:"devices"`
	Color   api.MicroColor `json:"color"`
}

// apiDevices
//   - GET - "/api/devices"
//   - POST - "/api/devices/color" - { devices: Device[]; color: number[] }
//   - TODO: Add color cache somehow
func apiDevices(e *echo.Echo, o Api) {
	e.GET(o.ServerPathPrefix+"/api/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.GET(o.ServerPathPrefix+"/api/devices", func(c echo.Context) error {
		return c.JSON(http.StatusOK, api.GetDevices(o.Config))
	})

	e.POST(o.ServerPathPrefix+"/api/devices/color", func(c echo.Context) error {
		var data RequestDevicesColorData
		err := json.NewDecoder(c.Request().Body).Decode(&data)
		if err != nil {
			return err
		}

		data.Devices = api.PostDevicesColor(o.Config, data.Color, data.Devices...)
		for di, dd := range data.Devices {
			for _, fd := range cache.Devices {
				if dd.Server.Addr != fd.Server.Addr {
					continue
				}

				// Only merge things changed after PostDevicesColor call
				fd.Color = dd.Color
				fd.Error = dd.Error
				fd.Online = dd.Online

				// Data to return
				data.Devices[di] = fd
			}
		}

		return c.JSON(http.StatusOK, data.Devices)
	})
}
