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
			device, err := getDeviceFromParamID(c, r)
			if err != nil {
				return err
			}

			// TODO: ...

			return echo.NewHTTPError(501, "Not Implemented")
		}
	}

	return nil
}

func APISetDeviceWhite(r *services.Registry, method string) echo.HandlerFunc {
	return func(c echo.Context) error {
		device, err := getDeviceFromParamID(c, r)
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
		return nil, echo.NewHTTPError(http.StatusBadRequest, err)
	}

	device, err := r.Device.Get(id)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("device not found: %v", err))
	}

	return device, nil
}
