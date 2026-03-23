package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/knackwurstking/picow-led/internal/components"
	"github.com/knackwurstking/picow-led/internal/components/dialogs"
	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/internal/services"
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
				Devices: formData.SelectedDevices, // TODO: Get selected devices from form data
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
