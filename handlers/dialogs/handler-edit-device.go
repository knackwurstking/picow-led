package dialogs

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/knackwurstking/picow-led/handlers/dialogs/components"
	"github.com/knackwurstking/picow-led/models"
	"github.com/knackwurstking/picow-led/services"
	"github.com/knackwurstking/picow-led/utils"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetEditDevice(c echo.Context) error {
	deviceID, err := utils.QueryParamDeviceID(c, "id", true)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.WrapError(err, "get device ID from query parameter"))
	}

	slog.Info("Get dialog for edit or create a device", "id", deviceID)

	var device *models.Device
	if deviceID > 0 {
		device, err = h.registry.Devices.Get(deviceID)
		if err != nil {
			if services.IsNotFoundError(err) {
				return echo.NewHTTPError(http.StatusNotFound,
					utils.WrapError(fmt.Errorf("device with ID %d not found", deviceID), "fetch device"))
			}

			return echo.NewHTTPError(http.StatusInternalServerError, utils.WrapError(err, "fetch device with ID %d", deviceID))
		}
	}

	if device != nil {
		slog.Info("Device found, rendering edit dialog")
		if err = components.EditDeviceDialog(device, false, nil).Render(
			c.Request().Context(), c.Response(),
		); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, utils.WrapError(err, "render edit device dialog"))
		}
	} else {
		slog.Info("Device not found, rendering new device dialog")
		if err = components.NewDeviceDialog(false, nil).Render(
			c.Request().Context(), c.Response(),
		); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, utils.WrapError(err, "render new device dialog"))
		}
	}

	return nil
}

func (h *Handler) PostEditDevice(c echo.Context) error {
	device, err := h.parseEditDeviceForm(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.WrapError(err, "parse device form"))
	}

	if !device.Validate() {
		validationError := fmt.Errorf("device validation, invalid form data %#v", device)

		if err := components.NewDeviceDialog(true, validationError).Render(
			c.Request().Context(), c.Response(),
		); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, utils.WrapError(err, "render device dialog"))
		}

		return echo.NewHTTPError(http.StatusBadRequest, validationError)
	}

	slog.Info("Submit handler from the create device dialog",
		"name", device.Name, "addr", device.Addr)

	if _, err := h.registry.Devices.Add(device); err != nil {
		databaseError := utils.WrapError(err, "add device %s", device.Name)

		if err := components.NewDeviceDialog(true, databaseError).Render(
			c.Request().Context(), c.Response(),
		); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, utils.WrapError(err, "render device dialog"))
		}

		return echo.NewHTTPError(http.StatusInternalServerError, databaseError)
	}

	c.Response().Header().Set("HX-Trigger", "reloadDevices")
	return nil
}

func (h *Handler) PutEditDevice(c echo.Context) error {
	deviceID, err := utils.QueryParamDeviceID(c, "id", true)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.WrapError(err, "get device ID from query parameter"))
	}

	device, err := h.parseEditDeviceForm(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.WrapError(err, "parse device form"))
	}
	device.ID = deviceID

	if !device.Validate() {
		validationError := fmt.Errorf("device validation, invalid form data %#v", device)

		if err := components.EditDeviceDialog(device, true, validationError).Render(
			c.Request().Context(), c.Response(),
		); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, utils.WrapError(err, "render device dialog"))
		}

		return echo.NewHTTPError(http.StatusBadRequest, validationError)
	}

	slog.Info("Submit handler from the edit device dialog",
		"id", device.ID, "name", device.Name, "addr", device.Addr)

	if err := h.registry.Devices.Update(device); err != nil {
		databaseError := utils.WrapError(err, "update device %s", device.Name)

		if err := components.EditDeviceDialog(device, true, databaseError).Render(
			c.Request().Context(), c.Response(),
		); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, utils.WrapError(err, "render device dialog"))
		}

		return echo.NewHTTPError(http.StatusInternalServerError, databaseError)
	}

	c.Response().Header().Set("HX-Trigger", "reloadDevices")
	return nil
}

func (h *Handler) parseEditDeviceForm(c echo.Context) (*models.Device, error) {
	host := c.FormValue("host")
	port := c.FormValue("port")
	name := c.FormValue("device-name")

	device, err := models.NewDevice(models.Addr(fmt.Sprintf("%s:%s", host, port)), name)
	if err != nil {
		return nil, utils.WrapError(err, "invalid device address: %s:%s", host, port)
	}

	return device, nil
}
