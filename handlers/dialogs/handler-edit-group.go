package dialogs

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/knackwurstking/picow-led/handlers/dialogs/components"
	"github.com/knackwurstking/picow-led/models"
	"github.com/knackwurstking/picow-led/services"
	"github.com/knackwurstking/picow-led/utils"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetEditGroup(c echo.Context) error {
	groupID, err := utils.QueryParamGroupID(c, "id", true)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.WrapError(err, "get group ID from query parameter"))
	}

	var group *models.Group
	if groupID > 0 {
		group, err = h.registry.Groups.Get(groupID)
		if err != nil {
			if services.IsNotFoundError(err) {
				return echo.NewHTTPError(http.StatusNotFound, utils.WrapError(fmt.Errorf("group with ID %d not found", groupID), "get group"))
			}
			return echo.NewHTTPError(http.StatusInternalServerError, utils.WrapError(err, "get group with ID %d", groupID))
		}
	}

	devices, err := h.registry.Devices.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.WrapError(err, "list devices"))
	}

	if group != nil {
		d := components.EditGroupDialog(group, devices, false, nil)
		if err := d.Render(c.Request().Context(), c.Response()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, utils.WrapError(err, "render edit group dialog"))
		}
	} else {
		d := components.NewGroupDialog(devices, nil, false, nil)
		if err := d.Render(c.Request().Context(), c.Response()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, utils.WrapError(err, "render new group dialog"))
		}
	}

	return nil
}

func (h *Handler) PostEditGroup(c echo.Context) error {
	group, err := h.parseGroupForm(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.WrapError(err, "parse group form values"))
	}

	if !group.Validate() {
		const message = "validate group"
		h.reRenderGroupDialogWithError(c, group, fmt.Errorf(message))
		return echo.NewHTTPError(http.StatusBadRequest, utils.WrapError(fmt.Errorf(message), "validate group"))
	}

	if _, err = h.registry.Groups.Add(group); err != nil {
		err = utils.WrapError(err, "add group")
		h.reRenderGroupDialogWithError(c, group, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.Response().Header().Set("HX-Trigger", "reloadGroups")
	return nil
}

func (h *Handler) PutEditGroup(c echo.Context) error {
	group, err := h.parseGroupForm(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.WrapError(err, "parse group form values"))
	}
	group.ID, err = utils.QueryParamGroupID(c, "id", false)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.WrapError(err, "get group ID"))
	}

	if !group.Validate() {
		const message = "validate group"
		h.reRenderGroupDialogWithError(c, group, fmt.Errorf(message))
		return echo.NewHTTPError(http.StatusBadRequest, utils.WrapError(fmt.Errorf(message), "validate group"))
	}

	if err = h.registry.Groups.Update(group); err != nil {
		err = utils.WrapError(err, "update group")
		h.reRenderGroupDialogWithError(c, group, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.Response().Header().Set("HX-Trigger", "reloadGroups")
	return nil
}

func (h *Handler) parseGroupForm(c echo.Context) (*models.Group, error) {
	formValues, err := c.FormParams()
	if err != nil {
		return nil, utils.WrapError(err, "get form parameters")
	}

	groupName := formValues.Get("group-name")
	if groupName == "" {
		return nil, utils.WrapError(fmt.Errorf("group name is required"), "parse group form: missing group name")
	}

	var deviceIDs []models.DeviceID
	for _, value := range formValues["devices"] {
		deviceID, err := strconv.Atoi(value)
		if err != nil {
			return nil, utils.WrapError(err, "invalid device ID in form")
		}
		deviceIDs = append(deviceIDs, models.DeviceID(deviceID))
	}

	return models.NewGroup(groupName, deviceIDs), nil
}

func (h *Handler) reRenderGroupDialogWithError(c echo.Context, group *models.Group, err error) {
	devices, err := h.registry.Devices.List()
	if err != nil {
		slog.Warn("list devices for edit group dialog with error", "error", err)
		return
	}

	if group.ID == 0 {
		d := components.NewGroupDialog(devices, group.Devices, true, err)
		if err := d.Render(c.Request().Context(), c.Response()); err != nil {
			slog.Warn("render edit group dialog with error", "error", err)
			return
		}
	} else {
		d := components.EditGroupDialog(group, devices, true, err)
		if err := d.Render(c.Request().Context(), c.Response()); err != nil {
			slog.Warn("render edit group dialog with error", "error", err)
			return
		}
	}
}
