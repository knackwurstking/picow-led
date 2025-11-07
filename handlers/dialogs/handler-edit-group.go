package dialogs

import (
	"fmt"
	"net/http"
	"slices"

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

		components.NewGroupDialog(preselectedDeviceIDs, devices, false, nil)
	}

	return nil
}

func (h *Handler) PostEditGroup(c echo.Context) error {
	// TODO: ...

	c.Response().Header().Set("HX-Trigger", "reloadGroups")
	return nil
}
func (h *Handler) PutEditGroup(c echo.Context) error {
	// TODO: ...

	c.Response().Header().Set("HX-Trigger", "reloadGroups")
	return nil
}
