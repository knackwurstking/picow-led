package handlers

import (
	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/labstack/echo/v4"
)

func APISetDeviceColor(r *services.Registry, method string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return echo.NewHTTPError(501, "Not Implemented")
	}
}

func APISetDeviceWhite(r *services.Registry, method string) echo.HandlerFunc {
	return func(c echo.Context) error {
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
