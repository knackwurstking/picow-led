package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	_api "github.com/knackwurstking/picow-led-server/pkg/api"
	_clients "github.com/knackwurstking/picow-led-server/pkg/clients"
)

func setupEndpointApiDevices(g *echo.Group) {
	g.GET("/devices", func(c echo.Context) error {
		return c.JSON(http.StatusOK, api.Devices)
	})

	g.GET("/device", func(c echo.Context) error {
		rD := _api.Device{}
		if status, err := readBodyData(c, &rD); status != http.StatusOK {
			return c.JSON(status, err.Error())
		}

		for _, d := range api.Devices {
			if d.Server.Addr == rD.Server.Addr {
				return c.JSON(http.StatusOK, d)
			}
		}

		return c.JSON(http.StatusBadRequest, fmt.Sprintf("device \"%s\" not found", rD.Server.Addr))
	})

	g.POST("/device", func(c echo.Context) error {
		rD := _api.Device{}
		if status, err := readBodyData(c, &rD); status != http.StatusOK {
			return c.JSON(status, err.Error())
		}

		err := api.Devices.Add(&rD)
		go clients.Emit(_clients.EventTypeDevices, api.Devices)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		saveConfig()
		return c.JSON(http.StatusOK, nil)
	})

	g.PUT("/device", func(c echo.Context) error {
		rD := _api.Device{}
		if status, err := readBodyData(c, &rD); status != http.StatusOK {
			return c.JSON(status, err.Error())
		}

		err := api.Devices.Update(&rD)
		go clients.Emit(_clients.EventTypeDevices, api.Devices)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		saveConfig()
		return c.JSON(http.StatusOK, nil)
	})

	g.DELETE("/device", func(c echo.Context) error {
		rD := _api.Device{}
		if status, err := readBodyData(c, &rD); status != http.StatusOK {
			return c.JSON(status, err.Error())
		}

		api.Devices.Remove(rD.Server.Addr)
		go clients.Emit(_clients.EventTypeDevices, api.Devices)

		saveConfig()
		return c.JSON(http.StatusOK, nil)
	})

	g.GET("/device/pins", func(c echo.Context) error {
		rD := _api.Device{}
		if status, err := readBodyData(c, &rD); status != http.StatusOK {
			return c.JSON(status, err.Error())
		}

		for _, d := range api.Devices {
			if d.Server.Addr == rD.Server.Addr {
				return c.JSON(http.StatusOK, d.Pins)
			}
		}

		return c.JSON(
			http.StatusBadRequest,
			fmt.Sprintf("device \"%s\" not found", rD.Server.Addr),
		)
	})

	g.POST("/device/pins", func(c echo.Context) error {
		rD := _api.Device{}
		if status, err := readBodyData(c, &rD); status != http.StatusOK {
			return c.JSON(status, err.Error())
		}

		for _, d := range api.Devices {
			if d.Server.Addr == rD.Server.Addr {
				err := d.SetPins(rD.Pins)
				go clients.Emit(_clients.EventTypeDevice, d)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, err.Error())
				}

				saveConfig()
				return c.JSON(http.StatusOK, nil)
			}
		}

		return c.JSON(
			http.StatusBadRequest,
			fmt.Sprintf("device \"%s\" not found", rD.Server.Addr),
		)
	})

	g.GET("/device/color", func(c echo.Context) error {
		rD := _api.Device{}
		if status, err := readBodyData(c, &rD); status != http.StatusOK {
			return c.JSON(status, err.Error())
		}

		for _, d := range api.Devices {
			if d.Server.Addr == rD.Server.Addr {
				return c.JSON(http.StatusOK, d.Color)
			}
		}

		return c.JSON(
			http.StatusBadRequest,
			fmt.Sprintf("device \"%s\" not found", rD.Server.Addr),
		)
	})

	g.POST("/device/color", func(c echo.Context) error {
		rD := _api.Device{}
		if status, err := readBodyData(c, &rD); status != http.StatusOK {
			return c.JSON(status, err.Error())
		}

		for _, d := range api.Devices {
			if d.Server.Addr == rD.Server.Addr {
				err := d.SetColor(rD.Color)
				go clients.Emit(_clients.EventTypeDevice, d)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, err.Error())
				}

				return c.JSON(http.StatusOK, nil)
			}
		}

		return c.JSON(
			http.StatusBadRequest,
			fmt.Sprintf("device \"%s\" not found", rD.Server.Addr),
		)
	})
}
