package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/knackwurstking/picow-led/internal/utils"
	"github.com/knackwurstking/picow-led/pkg/models"
	"github.com/labstack/echo/v4"
)

func APISetDeviceColor(r *services.Registry, method string) echo.HandlerFunc {
	switch method {
	case http.MethodPost:
		return func(c echo.Context) error {
			device, eerr := getDeviceFromParamID(c, r)
			if eerr != nil {
				return eerr
			}

			if !strings.Contains(string(device.Type), "RGB") {
				return echo.NewHTTPError(http.StatusBadRequest,
					fmt.Errorf("device does not support RGBW: %v", device.Type))
			}

			color := device.ToColor()

			if qColor, err := utils.ParseQueryColor(c); err != nil && err != utils.ErrNotFound {
				return echo.NewHTTPError(http.StatusBadRequest,
					fmt.Errorf("%v: color=%s", err, c.QueryParam("color")))
			} else if err == nil {
				color.Color[0] = qColor[0]
				color.Color[1] = qColor[1]
				color.Color[2] = qColor[2]
			}

			if err := r.Device.SetCurrentColor(device.ID, color.GetDuty(device.Type)); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError,
					fmt.Errorf("set device color: %v", err))
			}

			return nil
		}
	}

	return nil
}

func APISetDeviceRGBW(r *services.Registry, method string) echo.HandlerFunc {
	switch method {
	case http.MethodPost:
		return func(c echo.Context) error {
			device, eerr := getDeviceFromParamID(c, r)
			if eerr != nil {
				return eerr
			}

			if !strings.Contains(string(device.Type), "RGB") && !strings.Contains(string(device.Type), "W") {
				return echo.NewHTTPError(http.StatusBadRequest,
					fmt.Errorf("device does not support RGBW: %v", device.Type))
			}

			color := device.ToColor()

			if qColor, err := utils.ParseQueryColor(c); err != nil && err != utils.ErrNotFound {
				return echo.NewHTTPError(http.StatusBadRequest,
					fmt.Errorf("%v: color=%s", err, c.QueryParam("color")))
			} else if err == nil {
				color.Color[0] = qColor[0]
				color.Color[1] = qColor[1]
				color.Color[2] = qColor[2]
			}

			if qWhite, err := utils.ParseQueryWhite(c); err != nil && err != utils.ErrNotFound {
				return echo.NewHTTPError(http.StatusBadRequest,
					fmt.Errorf("%v white=%s", err, c.QueryParam("white")))
			} else if err == nil {
				color.White = qWhite
			}

			if err := r.Device.SetCurrentColor(device.ID, color.GetDuty(device.Type)); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError,
					fmt.Errorf("set device color: %v", err))
			}

			return nil
		}
	}

	return nil
}

func APISetDeviceWhite(r *services.Registry, method string) echo.HandlerFunc {
	switch method {
	case http.MethodPost:
		return func(c echo.Context) error {
			device, eerr := getDeviceFromParamID(c, r)
			if eerr != nil {
				return eerr
			}

			if strings.Contains(string(device.Type), "W") {
				return echo.NewHTTPError(http.StatusBadRequest,
					fmt.Errorf("device does not support white: %v", device.Type))
			}

			color := device.ToColor()

			if white, err := utils.ParseQueryWhite(c); err != nil && err != utils.ErrNotFound {
				return echo.NewHTTPError(http.StatusBadRequest,
					fmt.Errorf("%v: white: %s", err, c.QueryParam("white")))
			} else if err == nil {
				color.White = white
			}

			if err := r.Device.SetCurrentColor(device.ID, color.GetDuty(device.Type)); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError,
					fmt.Errorf("set device color: %v", err))
			}

			return nil
		}
	}

	return nil
}

func APISetDeviceWhite2(r *services.Registry, method string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return echo.NewHTTPError(501, "Not Implemented")
	}
}

func APISetDeviceBrightness(r *services.Registry, method string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return echo.NewHTTPError(501, "Not Implemented")
	}
}

func getDeviceFromParamID(c echo.Context, r *services.Registry) (*models.Device, *echo.HTTPError) {
	id, err := utils.ParseParamID(c)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest,
			fmt.Errorf("%v: id=%v", err, c.Param("id")))
	}

	device, err := r.Device.Get(id)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound,
			fmt.Errorf("device not found: %v", err))
	}

	return device, nil
}
