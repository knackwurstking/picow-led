package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/knackwurstking/picow-led/internal/templates/components/alert"
	"github.com/knackwurstking/picow-led/internal/templates/components/dialogs"
	"github.com/knackwurstking/picow-led/internal/templates/home"
	"github.com/knackwurstking/picow-led/internal/utils"
	"github.com/knackwurstking/picow-led/pkg/models"
	"github.com/labstack/echo/v4"
)

func HTMXGroups(r *services.Registry) echo.HandlerFunc {
	return func(c echo.Context) error {
		var errs []error

		groups, err := r.Group.List()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to list groups: %w", err))
		}

		t := home.Groups(home.GroupsProps{
			Data:   groups,
			Errors: errs,
		})

		if err := t.Render(c.Request().Context(), c.Response()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError,
				fmt.Errorf("failed to render template: %w", err))
		}

		return nil
	}
}

func HTMXPowerGroup(r *services.Registry) echo.HandlerFunc {
	return func(c echo.Context) error {
		mode := c.QueryParam("mode")

		groupIDStr := c.QueryParam("id")
		groupID, err := strconv.Atoi(groupIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid group ID: %v", err))
		}

		var errs []error

		if group, err := r.Group.Get(models.ID(groupID)); err != nil {
			errs = append(errs, fmt.Errorf("failed to get group: %w", err))
		} else {
			wg := &sync.WaitGroup{}
			for _, id := range group.Devices {
				wg.Go(func() {
					switch mode {
					case "on":
						if err := r.Device.TurnOn(id); err != nil {
							errs = append(errs, fmt.Errorf("failed to toggle power for device %d: %w", id, err))
						}
					case "off":
						if err := r.Device.TurnOff(id); err != nil {
							errs = append(errs, fmt.Errorf("failed to toggle power for device %d: %w", id, err))
						}
					default:
						errs = append(errs, fmt.Errorf("invalid mode: %s", mode))
					}
				})
			}
			wg.Wait()
		}

		// Handle errors (e.g. Render error messages)
		for _, err := range errs {
			if err := alert.RenderError(c, err.Error()); err != nil {
				return err
			}
		}

		c.Response().Header().Set("HX-Trigger", "reload-devices")

		return nil
	}
}

func HTMXAddGroupDialog(r *services.Registry, method string) echo.HandlerFunc {
	log := env.NewLogger("handlers.HTMXAddGroupDialog")

	parse := func(c echo.Context) (formData dialogs.AddGroupFormData, errs []error) {
		formData.Name = strings.TrimSpace(c.FormValue("name"))

		deviceIDs := strings.SplitSeq(c.FormValue("devices"), ",")
		for idString := range deviceIDs {
			id, _ := strconv.Atoi(idString)
			if d, err := r.Device.Get(models.ID(id)); err != nil {
				errs = append(errs, fmt.Errorf("failed to get device with ID %d: %w", id, err))
				continue
			} else {
				formData.SelectedDevices = append(formData.SelectedDevices, d.ID)
			}
		}

		return
	}

	render := func(c echo.Context, open bool, formData dialogs.AddGroupFormData, errs ...error) error {
		formData.Devices, _ = r.Device.List()
		t := dialogs.AddGroup(dialogs.AddGroupProps{
			AddGroupFormData: formData,
			Open:             open,
			OOB:              true,
			Errors:           errs,
		})
		if err := t.Render(c.Request().Context(), c.Response()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to render template").SetInternal(err)
		}
		return nil
	}

	switch method {
	case http.MethodGet:
		return func(c echo.Context) error {
			return render(c, true, dialogs.AddGroupFormData{})
		}
	case http.MethodPost:
		return func(c echo.Context) error {
			formData, errs := parse(c)
			if len(errs) > 0 {
				log.Error("failed to parse form data: %v", errs)
				return render(c, true, formData, errs...)
			}

			group := &models.Group{
				Name:    formData.Name,
				Devices: formData.SelectedDevices,
			}

			// Add group to registry group database
			if _, err := r.Group.Add(group); err != nil {
				log.Error("failed to add group: %v", err)
				errs = append(errs, fmt.Errorf("failed to add group: %w", err))
				return render(c, true, formData, errs...)
			}

			c.Response().Header().Set("HX-Trigger", "reload-groups")

			return render(c, false, formData, errs...) // Close dialog on success
		}
	}

	return nil
}

func HTMXEditGroupDialog(r *services.Registry, method string) echo.HandlerFunc {
	log := env.NewLogger("handlers.HTMXEditGroupDialog")

	parseForm := func(c echo.Context) (formData dialogs.EditGroupFormData, errs []error) {
		id, err := utils.ParseQueryID(c)
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid group ID: %w", err))
		}
		formData.ID = id

		formData.Name = strings.TrimSpace(c.FormValue("name"))

		deviceIDs := strings.SplitSeq(c.FormValue("devices"), ",")
		for idString := range deviceIDs {
			id, _ := strconv.Atoi(idString)
			if d, err := r.Device.Get(models.ID(id)); err != nil {
				errs = append(errs, fmt.Errorf("failed to get device with ID %d: %w", id, err))
				continue
			} else {
				formData.SelectedDevices = append(formData.SelectedDevices, d.ID)
			}
		}

		return
	}

	render := func(c echo.Context, open bool, formData dialogs.EditGroupFormData, errs ...error) error {
		formData.Devices, _ = r.Device.List()
		t := dialogs.EditGroup(dialogs.EditGroupProps{
			EditGroupFormData: formData,
			Open:              open,
			OOB:               true,
			Errors:            errs,
		})
		if err := t.Render(c.Request().Context(), c.Response()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to render template").SetInternal(err)
		}
		return nil
	}

	switch method {
	case http.MethodGet:
		return func(c echo.Context) error {
			var errs []error

			id, err := utils.ParseQueryID(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid group ID: %v", err))
			}

			formData := dialogs.EditGroupFormData{
				ID: id,
			}

			if group, err := r.Group.Get(models.ID(id)); err != nil {
				errs = append(errs, fmt.Errorf("failed to get group: %w", err))
			} else {
				formData.Name = group.Name
				formData.SelectedDevices = group.Devices
			}

			return render(c, true, formData, errs...)
		}
	case http.MethodPost:
		return func(c echo.Context) error {
			formData, errs := parseForm(c)
			if len(errs) > 0 {
				return render(c, true, formData, errs...)
			}

			group := &models.Group{
				ID:      formData.ID,
				Name:    formData.Name,
				Devices: formData.SelectedDevices,
			}

			if err := r.Group.Update(group); err != nil {
				log.Error("failed to update group: %v", err)
				errs = append(errs, fmt.Errorf("failed to update group: %w", err))
				return render(c, true, formData, errs...)
			}

			c.Response().Header().Set("HX-Trigger", "reload-groups")

			return render(c, false, formData, errs...)
		}
	case http.MethodDelete:
		return func(c echo.Context) error {
			id, err := utils.ParseQueryID(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest,
					fmt.Errorf("invalid group ID: %v", err))
			}

			if err := r.Group.Delete(id); err != nil {
				return alert.RenderError(c,
					fmt.Sprintf("Failed to delete group with ID %d: %v", id, err))
			}

			c.Response().Header().Set("HX-Trigger", "reload-groups")

			return render(c, false, dialogs.EditGroupFormData{}) // Close dialog
		}
	}

	return nil
}
