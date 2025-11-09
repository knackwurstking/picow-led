package home

import (
	"log/slog"

	"github.com/knackwurstking/picow-led/handlers/home/components"
	"github.com/knackwurstking/picow-led/models"
	"github.com/labstack/echo/v4"
)

// ******* //
// Devices //
// ******* //

func OOBRenderPageHomeDeviceError(c echo.Context, deviceID models.DeviceID, err error) {
	deviceError := components.OOBDeviceError(deviceID, err, true)
	if err := deviceError.Render(c.Request().Context(), c.Response()); err != nil {
		slog.Error("Failed to render device error page", "deviceID", deviceID, "error", err)
	}
}

func OOBRenderPageHomeDevicePowerButton(c echo.Context, deviceID models.DeviceID, currentColor []uint8) {
	devicePowerButton := components.OOBDevicePowerButton(deviceID, currentColor, true)
	if err := devicePowerButton.Render(c.Request().Context(), c.Response()); err != nil {
		slog.Error("Failed to render device power button page", "deviceID", deviceID, "error", err)
	}
}

// ****** //
// Groups //
// ****** //

func OOBRenderPageHomeGroupError(c echo.Context, groupID models.GroupID, err []error) {
	deviceError := components.OOBGroupError(groupID, err, true)
	if err := deviceError.Render(c.Request().Context(), c.Response()); err != nil {
		slog.Error("Failed to render device error page", "deviceID", groupID, "error", err)
	}
}
