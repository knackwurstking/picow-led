package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/knackwurstking/picow-led/internal/components"
	"github.com/knackwurstking/picow-led/internal/components/dialogs"
	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/internal/services"
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

		t := components.Groups(components.GroupsProps{
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
			t := components.OOBAddAlert(env.IDAlertContainer, components.AlertTypeError, err.Error())
			if err = t.Render(c.Request().Context(), c.Response()); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to render error alert: %w", err))
			}
		}

		return nil
	}
}

func HTMXAddGroupDialog(r *services.Registry, method string) echo.HandlerFunc {
	log := env.NewLogger("handlers.HTMXAddGroupDialog")

	parse := func(c echo.Context) (formData dialogs.AddGroupFormData, errs []error) {
		formData.Name = strings.TrimSpace(c.FormValue("name"))

		deviceIDs := strings.Split(c.FormValue("devices"), ",")
		for _, idString := range deviceIDs {
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
			var errs []error

			devices, err := r.Device.List()
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to list devices: %w", err))
			}

			return render(c, true, dialogs.AddGroupFormData{
				Devices: devices,
			}, errs...)
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

			return render(c, false, formData, errs...) // Close dialog on success
		}
	}

	return nil
}

func HTMXEditGroupDialog(r *services.Registry, method string) echo.HandlerFunc {
	parseForm := func(c echo.Context) (dialogs.EditGroupFormData, []error) {
		var errs []error
		var formData dialogs.EditGroupFormData

		// TODO: ...

		return formData, errs
	}

	render := func(c echo.Context, open bool, formData dialogs.EditGroupFormData, errs ...error) error {
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
				if devices, err := r.Device.List(); err != nil {
					errs = append(errs, fmt.Errorf("failed to list devices: %w", err))
				} else {
					formData.Name = group.Name
					formData.Devices = devices
					formData.SelectedDevices = group.Devices
				}
			}

			return render(c, true, formData, errs...)
		}
	case http.MethodPost:
		return func(c echo.Context) error {
			formData, errs := parseForm(c)
			if len(errs) > 0 {
				return render(c, true, formData, errs...)
			}

			// TODO: ...

			return echo.NewHTTPError(http.StatusNotImplemented, "not implemented")
		}
	case http.MethodDelete:
		return func(c echo.Context) error {
			// TODO: ...

			return echo.NewHTTPError(http.StatusNotImplemented, "not implemented")
		}
	}

	return nil
}
