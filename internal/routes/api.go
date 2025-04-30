package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

// apiRoutes
//   - apiSetupPing: 	GET 	- "/api/ping"
//   - apiSetupDevices: GET 	- "/api/devices"
//   - apiSetupDevices: POST 	- "/api/devices/color" <- { devices: Device[]; color: number[] }
//   - apiSetupColor: 	GET 	- "/api/color"
//   - apiSetupColor: 	GET 	- "/api/color/:name"
//   - TODO: apiSetupColor: 	POST 	- "/api/color/:name" <- `number[]`
func apiRoutes(e *echo.Echo, o Api) {
	apiSetupPing(e, o)
	apiSetupDevices(e, o)
	apiSetupColor(e, o)
}

func apiSetupPing(e *echo.Echo, o Api) {
	e.GET(o.ServerPathPrefix+"/api/ping", func(c echo.Context) error {
		err := c.String(http.StatusOK, "pong")
		if err != nil {
			log.Println(err)
		}
		return err
	})
}

func apiSetupDevices(e *echo.Echo, o Api) {
	e.GET(o.ServerPathPrefix+"/api/devices", func(c echo.Context) error {
		err := c.JSON(http.StatusOK, api.GetDevices(o.Config))
		if err != nil {
			log.Println(err)
		}
		return err
	})

	e.POST(o.ServerPathPrefix+"/api/devices/color", func(c echo.Context) error {
		var data RequestDevicesColorData
		err := json.NewDecoder(c.Request().Body).Decode(&data)
		if err != nil {
			log.Println(err)
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

		err = c.JSON(http.StatusOK, data.Devices)
		if err != nil {
			log.Println(err)
		}
		return err
	})
}

func apiSetupColor(e *echo.Echo, o Api) {
	e.GET(o.ServerPathPrefix+"/api/color", func(c echo.Context) error {
		err := c.JSON(http.StatusOK, cache.Color)
		if err != nil {
			log.Println(err)
		}
		return err
	})

	e.GET(o.ServerPathPrefix+"/api/color/:name", func(c echo.Context) error {
		name := url.QueryEscape(c.Param("name"))
		color, ok := cache.Color[name]
		if !ok {
			err := c.String(http.StatusBadRequest, fmt.Sprintf("Color for \"%s\" not found", name))
			if err != nil {
				log.Println(err)
			}
			return err
		}

		err := c.JSON(http.StatusOK, color)
		if err != nil {
			log.Println(err)
		}
		return err
	})
}
