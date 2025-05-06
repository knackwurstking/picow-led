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

type apiOptions struct {
	ServerPathPrefix string
	Config           *api.Config
}

func apiRoutes(e *echo.Echo, o apiOptions) {
	apiSetupPing(e, o)
	apiSetupDevices(e, o)
	apiSetupColor(e, o)
}

func apiSetupPing(e *echo.Echo, o apiOptions) {
	e.GET(o.ServerPathPrefix+"/api/ping", func(c echo.Context) error {
		err := c.String(http.StatusOK, "pong")
		if err != nil {
			log.Println(err)
		}
		return err
	})
}

func apiSetupDevices(e *echo.Echo, o apiOptions) {
	e.GET(o.ServerPathPrefix+"/api/devices", func(c echo.Context) error {
		devices := api.GetDevices(o.Config)
		err := c.JSON(http.StatusOK, devices)
		if err != nil {
			log.Println(err)
		}

		go cache.SetDevices(devices...)
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
			return c.String(http.StatusBadRequest, err.Error())
		}

		reqData.Devices = api.PostDevicesColor(o.Config, reqData.Color, reqData.Devices...)

		for i, d := range reqData.Devices {
			d, err := cache.UpdateDevice(d.Server.Addr, d)
			if err != nil {
				log.Println(err)
				return c.String(http.StatusNotFound, err.Error())
			}

			reqData.Devices[i] = d
		}

		err = c.JSON(http.StatusOK, reqData.Devices)
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return nil
	})
}

func apiSetupColor(e *echo.Echo, o apiOptions) {
	e.GET(o.ServerPathPrefix+"/api/color", func(c echo.Context) error {
		err := c.JSON(http.StatusOK, cache.Color())
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

		cacheColor := cache.Color()
		if len(cacheColor)-1 < index {
			return c.String(http.StatusBadRequest, fmt.Sprintf("color index %d not found", index))
		}
		color := cacheColor[index]

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

		var reqData api.MicroColor
		err = json.NewDecoder(c.Request().Body).Decode(&reqData)
		if err != nil {
			log.Println(err)
			return err
		}

		err = cache.UpdateColor(index, reqData)
		if err != nil {
			log.Println(err.Error())
			return c.String(http.StatusBadRequest, err.Error())
		}

		return nil
	})
}
