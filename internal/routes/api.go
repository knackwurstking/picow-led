// API Routes:
//   - apiSetupPing: 	GET 	- "/api/ping"
//   - apiSetupDevices: GET 	- "/api/devices"
//   - apiSetupDevices: POST 	- "/api/devices/color" <- { devices: Device[]; color: number[] }
//   - apiSetupColor: 	GET 	- "/api/color"
//   - apiSetupColor: 	GET 	- "/api/color/:index"
//   - apiSetupColor: 	POST 	- "/api/color/:index" <- `number[]`
package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"picow-led/internal/api"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Api struct {
	ServerPathPrefix string
	Config           *api.Config
}

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
		devices := api.GetDevices(o.Config)
		err := c.JSON(http.StatusOK, devices)
		if err != nil {
			log.Println(err)
		}

		go func() {
			cache.Mutex.Lock()
			defer cache.Mutex.Unlock()

			cache.Devices = devices
		}()

		return err
	})

	e.POST(o.ServerPathPrefix+"/api/devices/color", func(c echo.Context) error {
		var reqData struct {
			Devices []*api.Device  `json:"devices"`
			Color   api.MicroColor `json:"color"`
		}
		err := json.NewDecoder(c.Request().Body).Decode(&reqData)
		if err != nil {
			log.Println(err)
			return err
		}

		reqData.Devices = api.PostDevicesColor(o.Config, reqData.Color, reqData.Devices...)

		func() {
			cache.Mutex.Lock()
			defer cache.Mutex.Unlock()

			for di, dd := range reqData.Devices {
				for _, fd := range cache.Devices {
					if dd.Server.Addr != fd.Server.Addr {
						continue
					}

					// Only merge things changed after PostDevicesColor call
					fd.Color = dd.Color
					fd.Error = dd.Error
					fd.Online = dd.Online

					// Data to return
					reqData.Devices[di] = fd
				}
			}
		}()

		err = c.JSON(http.StatusOK, reqData.Devices)
		if err != nil {
			log.Println(err)
		}
		return err
	})
}

func apiSetupColor(e *echo.Echo, o Api) {
	e.GET(o.ServerPathPrefix+"/api/color", func(c echo.Context) error {
		cache.Mutex.Lock()
		defer cache.Mutex.Unlock()

		err := c.JSON(http.StatusOK, cache.Color)
		if err != nil {
			log.Println(err)
		}
		return err
	})

	e.GET(o.ServerPathPrefix+"/api/color/:index", func(c echo.Context) error {
		index, err := strconv.Atoi(c.Param("index"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		cache.Mutex.Lock()
		defer cache.Mutex.Unlock()

		if len(cache.Color)-1 < index {
			return c.String(http.StatusBadRequest, fmt.Sprintf("color index %d not found", index))
		}
		color := cache.Color[index]

		err = c.JSON(http.StatusOK, color)
		if err != nil {
			log.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		return nil
	})

	e.POST(o.ServerPathPrefix+"/api/color/:index", func(c echo.Context) error {
		index, err := strconv.Atoi(c.Param("index"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		cache.Mutex.Lock()
		defer cache.Mutex.Unlock()

		if len(cache.Color)-1 < index {
			return c.String(http.StatusBadRequest, fmt.Sprintf("index %d not exists", index))
		}

		var reqData api.MicroColor
		err = json.NewDecoder(c.Request().Body).Decode(&reqData)
		if err != nil {
			log.Println(err)
			return err
		}

		cache.Color[index] = reqData

		return nil
	})
}
