package dialogs

import (
	"fmt"
	"log/slog"
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
	group, err := h.parseGroupForm(c)
	if err != nil {
		return fmt.Errorf("failed to parse group form values: %v", err)
	}

	if !group.Validate() {
		message := "failed to validate group"
		h.reRenderGroupDialogWithError(c, group, fmt.Errorf(message))
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf(message))
	}

	if _, err = h.registry.Groups.Add(group); err != nil {
		err = fmt.Errorf("failed to add group: %v", err)
		h.reRenderGroupDialogWithError(c, group, err)
		return err
	}

	c.Response().Header().Set("HX-Trigger", "reloadGroups")
	return nil
}

func (h *Handler) PutEditGroup(c echo.Context) error {
	group, err := h.parseGroupForm(c)
	if err != nil {
		return fmt.Errorf("failed to parse group form values: %v", err)
	}
	group.ID, err = utils.QueryParamGroupID(c, "id", false)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to get group ID: %v", err))
	}

	if !group.Validate() {
		message := "failed to validate group"
		h.reRenderGroupDialogWithError(c, group, fmt.Errorf(message))
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf(message))
	}

	if err = h.registry.Groups.Update(group); err != nil {
		err = fmt.Errorf("failed to update group: %v", err)
		h.reRenderGroupDialogWithError(c, group, err)
		return err
	}

	c.Response().Header().Set("HX-Trigger", "reloadGroups")
	return nil
}

func (h *Handler) parseGroupForm(c echo.Context) (*models.Group, error) {
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

func (h *Handler) reRenderGroupDialogWithError(c echo.Context, group *models.Group, err error) {
	devices, err := h.registry.Devices.List()
	if err != nil {
		slog.Warn("Failed to list devices for edit group dialog with error", "error", err)
		return
	}

	if group.ID == 0 {
		d := components.NewGroupDialog(group.Devices, devices, true, err)
		if err := d.Render(c.Request().Context(), c.Response()); err != nil {
			slog.Warn("Failed to render edit group dialog with error", "error", err)
			return
		}
	} else {
		d := components.EditGroupDialog(group, devices, true, err)
		if err := d.Render(c.Request().Context(), c.Response()); err != nil {
			slog.Warn("Failed to render edit group dialog with error", "error", err)
			return
		}
	}
}
