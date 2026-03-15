package handlers

import (
	"fmt"
	"net/http"

	"github.com/knackwurstking/picow-led/internal/pages"
	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	t := pages.Home()
	if err := t.Render(c.Request().Context(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to render template: %v", err))
	}
	return nil
}
