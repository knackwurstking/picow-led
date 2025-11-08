package dialogs

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/knackwurstking/picow-led/handlers/dialogs/components"
	"github.com/knackwurstking/picow-led/handlers/utils"
	"github.com/knackwurstking/picow-led/models"
	"github.com/knackwurstking/picow-led/services"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetEditGroup(c echo.Context) error {
	groupID, err := utils.QueryParamGroupID(c, "id", true)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var group *models.Group
	if groupID > 0 {
		group, err = h.registry.Groups.Get(groupID)
		if err != nil {
			if services.IsNotFoundError(err) {
				return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("group with ID %d not found", groupID))
			}
			return fmt.Errorf("failed to get group with ID %d: %v", groupID, err)
		}
	}

	devices, err := h.registry.Devices.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if group != nil {
		components.EditGroupDialog(group, devices, false, nil)
	} else {
		var preselectedDeviceIDs []models.DeviceID

		// Preselect devices based on their current color state
		for _, d := range devices {
			color, _ := h.registry.DeviceControls.GetCurrentColor(d.ID)
			if len(color) > 0 && slices.Max(color) > 0 {
				preselectedDeviceIDs = append(preselectedDeviceIDs, d.ID)
			}
		}

		d := components.NewGroupDialog(preselectedDeviceIDs, devices, false, nil)
		if err := d.Render(c.Request().Context(), c.Response()); err != nil {
			return fmt.Errorf("failed to render group dialog: %v", err)
		}
	}

	return nil
}

func (h *Handler) PostEditGroup(c echo.Context) error {
	group, err := h.parseEditGroupForm(c)
	if err != nil {
		// TODO: re render the dialog (with oob set to true) and add the error to the dialog

		return fmt.Errorf("failed to parse group form values: %v", err)
	}

	// TODO: validate group, render dialog on error, add group to the database using `h.registry`, set the hx trigger and return

	c.Response().Header().Set("HX-Trigger", "reloadGroups")
	return nil
}

func (h *Handler) PutEditGroup(c echo.Context) error {

	// TODO: ...

	c.Response().Header().Set("HX-Trigger", "reloadGroups")
	return nil
}

func (h *Handler) parseEditGroupForm(c echo.Context) (*models.Group, error) {
	formValues, err := c.FormParams()
	if err != nil {
		return nil, fmt.Errorf("failed to get form parameters: %v", err)
	}

	groupName := formValues.Get("group-name")
	if groupName == "" {
		return nil, fmt.Errorf("group name is required")
	}

	var deviceIDs []models.DeviceID
	for _, value := range formValues["devices"] {
		deviceID, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid device ID: %v", err)
		}
		deviceIDs = append(deviceIDs, models.DeviceID(deviceID))
	}

	return models.NewGroup(groupName, deviceIDs), nil
}
