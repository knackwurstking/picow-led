package handlers

import (
	"net/http"

	"github.com/knackwurstking/picow-led/components"
	"github.com/knackwurstking/picow-led/services"
	"github.com/labstack/echo/v4"
)

type Pages struct {
	registry *services.Registry
}

func NewPages(r *services.Registry) *Pages {
	return &Pages{
		registry: r,
	}
}

func (p Pages) Register(e *echo.Echo) {
	Register(e, http.MethodGet, "", p.Home)
}

func (p *Pages) Home(c echo.Context) error {
	err := components.PageHome().Render(c.Request().Context(), c.Response())
	if err != nil {
		return err
	}

	return nil
}
