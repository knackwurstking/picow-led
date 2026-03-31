package handlers

import (
	"fmt"
	"net/http"

	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/knackwurstking/picow-led/internal/templates/components/alert"
	devicepage "github.com/knackwurstking/picow-led/internal/templates/device"
	"github.com/knackwurstking/picow-led/internal/templates/home"
	"github.com/knackwurstking/picow-led/internal/utils"
	"github.com/knackwurstking/picow-led/pkg/models"
	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	t := home.Page()
	if err := t.Render(c.Request().Context(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to render template: %v", err))
	}
	return nil
}

func Device(r *services.Registry, method string) echo.HandlerFunc {
	log := env.NewLogger("handlers.Device")

	render := func(c echo.Context, device *models.Device) error {
		pins, err := r.Device.GetPins(device.ID)
		if err != nil {
			defer func() {
				message := fmt.Sprintf("Failed to load device pins: %s", err)
				if err := alert.RenderError(c, message); err != nil {
					log.Error("Failed to render error alert: %v", err)
				}
			}()
		}

		t := devicepage.Page(device, pins...)
		if err := t.Render(c.Request().Context(), c.Response()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Failed to render page: %w", err))
		}
		return nil
	}

	switch method {
	case http.MethodGet:
		return func(c echo.Context) error {
			id, err := utils.ParseQueryID(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%v: %s", err, c.QueryParam("id")))
			}

			device, err := r.Device.Get(id)
			if err != nil {
				return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Device not found: %w", err))
			}

			return render(c, device)
		}
	}

	return nil
}
