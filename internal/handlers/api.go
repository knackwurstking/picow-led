package handlers

import (
	"fmt"
	"net/http"

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

			color, err := utils.ParseQueryColor(c)
			if err != nil && err != utils.ErrNotFound {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid color: %v", err))
			} else if err == utils.ErrNotFound {
				color = device.Color
			}

			if err := r.Device.SetCurrentColor(device.ID, color); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("set device color: %v", err))
			}

			return nil
		}
	}

	return nil
}

func APISetDeviceWhite(r *services.Registry, method string) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := getDeviceFromParamID(c, r)
		if err != nil {
			return err
		}

		// TODO: ...

		return echo.NewHTTPError(501, "Not Implemented")
	}
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
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%v: %v", err, c.Param("id")))
	}

	device, err := r.Device.Get(id)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("device not found: %v", err))
	}

	return device, nil
}
