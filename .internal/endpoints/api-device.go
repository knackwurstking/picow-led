package endpoints

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/knackwurstking/picow-led-server/pkg/api"
)

func createApiDeviceEndpoints(g *echo.Group, a *api.API, changeCallback func()) {
	g.GET("/device", func(ctx echo.Context) error {
		reqData := api.Device{}
		if status, err := readBodyData(ctx, &reqData); err != nil {
			return ctx.String(status, err.Error())
		}

		for _, d := range a.Devices {
			if d.Server.Addr == reqData.Server.Addr {
				return ctx.JSON(http.StatusOK, d)
			}
		}

		return ctx.String(
			http.StatusBadRequest,
			fmt.Sprintf(
				"device \"%s\" not found",
				reqData.Server.Addr,
			),
		)
	})

	g.POST("/device", func(ctx echo.Context) error {
		reqData := api.Device{}
		if status, err := readBodyData(ctx, &reqData); err != nil {
			return ctx.String(status, err.Error())
		}

		err := a.AddDevice(&reqData)
		if err != nil {
			return ctx.String(http.StatusBadRequest, err.Error())
		}

		defer changeCallback()
		return ctx.JSON(http.StatusOK, nil)
	})

	g.PUT("/device", func(ctx echo.Context) error {
		reqData := api.Device{}
		if status, err := readBodyData(ctx, &reqData); err != nil {
			return ctx.String(status, err.Error())
		}

		err := a.UpdateDevice(&reqData)
		if err != nil {
			return ctx.String(http.StatusBadRequest, err.Error())
		}

		defer changeCallback()
		return ctx.JSON(http.StatusOK, nil)
	})

	g.DELETE("/device", func(ctx echo.Context) error {
		reqData := api.Device{}
		if status, err := readBodyData(ctx, &reqData); err != nil {
			return ctx.String(status, err.Error())
		}

		a.DeleteDevice(reqData.Server.Addr)

		defer changeCallback()
		return ctx.JSON(http.StatusOK, nil)
	})

	g.GET("/device/pins", func(ctx echo.Context) error {
		reqData := api.Device{}
		if status, err := readBodyData(ctx, &reqData); err != nil {
			return ctx.String(status, err.Error())
		}

		for _, d := range a.Devices {
			if d.Server.Addr == reqData.Server.Addr {
				return ctx.JSON(http.StatusOK, d.Pins)
			}
		}

		return ctx.String(
			http.StatusBadRequest,
			fmt.Sprintf(
				"device \"%s\" not found",
				reqData.Server.Addr,
			),
		)
	})

	g.POST("/device/pins", func(ctx echo.Context) error {
		reqData := api.Device{}
		if status, err := readBodyData(ctx, &reqData); err != nil {
			return ctx.String(status, err.Error())
		}

		if err := a.UpdateDevicePins(reqData.Server.Addr, reqData.Pins); err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		} else {
			defer changeCallback()
			return ctx.JSON(http.StatusOK, nil)
		}
	})

	g.GET("/device/color", func(ctx echo.Context) error {
		reqData := api.Device{}
		if status, err := readBodyData(ctx, &reqData); err != nil {
			return ctx.String(status, err.Error())
		}

		for _, d := range a.Devices {
			if d.Server.Addr == reqData.Server.Addr {
				return ctx.JSON(http.StatusOK, d.Color)
			}
		}

		return ctx.String(
			http.StatusBadRequest,
			fmt.Sprintf("device \"%s\" not found", reqData.Server.Addr),
		)
	})

	g.POST("/device/color", func(ctx echo.Context) error {
		reqData := api.Device{}
		if status, err := readBodyData(ctx, &reqData); err != nil {
			return ctx.String(status, err.Error())
		}

		if err := a.UpdateDeviceColor(reqData.Server.Addr, reqData.Color); err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		} else {
			return ctx.JSON(http.StatusOK, nil)
		}
	})
}
