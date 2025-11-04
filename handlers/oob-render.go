package handlers

import (
	"log/slog"

	"github.com/knackwurstking/picow-led/components"
	"github.com/knackwurstking/picow-led/models"
	"github.com/labstack/echo/v4"
)

func OOBRenderPageHomeDeviceError(c echo.Context, deviceID models.DeviceID, err error) {
	deviceError := components.PageHome_SectionDevices_DeviceError(deviceID, err, true)
	if err := deviceError.Render(c.Request().Context(), c.Response()); err != nil {
		slog.Error("Failed to render device error page", "deviceID", deviceID, "error", err)
	}
}
