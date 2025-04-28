package routes

import (
	"encoding/json"
	"net/http"
	"picow-led/internal/api"

	"github.com/labstack/echo/v4"
)

type RequestDevicesColorData struct {
	Devices []*api.Device  `json:"devices"`
	Color   api.MicroColor `json:"color"`
}

// apiDevices
//   - GET - "/api/devices"
//   - POST - "/api/devices/color" - { devices: Device[]; color: number[] }
func apiDevices(e *echo.Echo, o Options) {
	e.GET(o.ServerPathPrefix+"/api/devices", func(c echo.Context) error {
		return c.JSON(http.StatusOK, api.GetDevices(o.Api))
	})

	e.POST(o.ServerPathPrefix+"/api/devices/color", func(c echo.Context) error {
		var data RequestDevicesColorData
		err := json.NewDecoder(c.Request().Body).Decode(&data)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK,
			api.PostDevicesColor(o.Api, data.Color, data.Devices...),
		)
	})
}
