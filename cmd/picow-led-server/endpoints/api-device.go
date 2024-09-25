package endpoints

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/knackwurstking/picow-led-server/pkg/api"
	"github.com/knackwurstking/picow-led-server/pkg/clients"
)

func createApiDeviceEndpoints(g *echo.Group, c *clients.Clients, a *api.API, apiChangeCallback func()) {
	g.GET("/device", func(ctx echo.Context) error {
		rD := api.Device{}
		if status, err := readBodyData(ctx, &rD); status != http.StatusOK {
			return ctx.String(status, err.Error())
		}

		for _, d := range a.Devices {
			if d.Server.Addr == rD.Server.Addr {
				return ctx.JSON(http.StatusOK, d)
			}
		}

		return ctx.String(http.StatusBadRequest, fmt.Sprintf("device \"%s\" not found", rD.Server.Addr))
	})

	g.POST("/device", func(ctx echo.Context) error {
		rD := api.Device{}
		if status, err := readBodyData(ctx, &rD); status != http.StatusOK {
			return ctx.String(status, err.Error())
		}

		err := a.Devices.Add(&rD)
		go c.Emit(clients.EventTypeDevices, a.Devices)
		if err != nil {
			return ctx.String(http.StatusBadRequest, err.Error())
		}

		defer apiChangeCallback()
		return ctx.JSON(http.StatusOK, nil)
	})

	g.PUT("/device", func(ctx echo.Context) error {
		rD := api.Device{}
		if status, err := readBodyData(ctx, &rD); status != http.StatusOK {
			return ctx.String(status, err.Error())
		}

		err := a.Devices.Update(&rD)
		go c.Emit(clients.EventTypeDevices, a.Devices)
		if err != nil {
			return ctx.String(http.StatusBadRequest, err.Error())
		}

		defer apiChangeCallback()
		return ctx.JSON(http.StatusOK, nil)
	})

	g.DELETE("/device", func(ctx echo.Context) error {
		rD := api.Device{}
		if status, err := readBodyData(ctx, &rD); status != http.StatusOK {
			return ctx.String(status, err.Error())
		}

		a.Devices.Remove(rD.Server.Addr)
		go c.Emit(clients.EventTypeDevices, a.Devices)

		defer apiChangeCallback()
		return ctx.JSON(http.StatusOK, nil)
	})

	g.GET("/device/pins", func(ctx echo.Context) error {
		rD := api.Device{}
		if status, err := readBodyData(ctx, &rD); status != http.StatusOK {
			return ctx.String(status, err.Error())
		}

		for _, d := range a.Devices {
			if d.Server.Addr == rD.Server.Addr {
				return ctx.JSON(http.StatusOK, d.Pins)
			}
		}

		return ctx.String(
			http.StatusBadRequest,
			fmt.Sprintf("device \"%s\" not found", rD.Server.Addr),
		)
	})

	g.POST("/device/pins", func(ctx echo.Context) error {

		rD := api.Device{}
		if status, err := readBodyData(ctx, &rD); status != http.StatusOK {
			return ctx.String(status, err.Error())
		}

		for _, d := range a.Devices {
			if d.Server.Addr == rD.Server.Addr {
				err := d.SetPins(rD.Pins)
				go c.Emit(clients.EventTypeDevice, d)
				if err != nil {
					return ctx.String(http.StatusInternalServerError, err.Error())
				}

				defer apiChangeCallback()
				return ctx.JSON(http.StatusOK, nil)
			}
		}

		return ctx.String(
			http.StatusBadRequest,
			fmt.Sprintf("device \"%s\" not found", rD.Server.Addr),
		)
	})

	g.GET("/device/color", func(ctx echo.Context) error {
		rD := api.Device{}
		if status, err := readBodyData(ctx, &rD); status != http.StatusOK {
			return ctx.String(status, err.Error())
		}

		for _, d := range a.Devices {
			if d.Server.Addr == rD.Server.Addr {
				return ctx.JSON(http.StatusOK, d.Color)
			}
		}

		return ctx.String(
			http.StatusBadRequest,
			fmt.Sprintf("device \"%s\" not found", rD.Server.Addr),
		)
	})

	g.POST("/device/color", func(ctx echo.Context) error {
		rD := api.Device{}
		if status, err := readBodyData(ctx, &rD); status != http.StatusOK {
			return ctx.String(status, err.Error())
		}

		for _, d := range a.Devices {
			if d.Server.Addr == rD.Server.Addr {
				err := d.SetColor(rD.Color)
				go c.Emit(clients.EventTypeDevice, d)
				if err != nil {
					return ctx.String(http.StatusInternalServerError, err.Error())
				}

				return ctx.JSON(http.StatusOK, nil)
			}
		}

		return ctx.String(
			http.StatusBadRequest,
			fmt.Sprintf("device \"%s\" not found", rD.Server.Addr),
		)
	})
}
