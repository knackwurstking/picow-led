package handlers

import (
	"fmt"

	"github.com/knackwurstking/picow-led/components"
	"github.com/knackwurstking/picow-led/models"
	"github.com/labstack/echo/v4"
)

func OOBRenderPageHomeDeviceError(c echo.Context, deviceID models.DeviceID, err error) {
	// TODO: Swap error with "#device-%d-error"
	id := fmt.Sprintf(components.IDPageHome_SectionDevices_DeviceError, deviceID)
}
