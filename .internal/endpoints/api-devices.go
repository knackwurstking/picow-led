package endpoints

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/knackwurstking/picow-led-server/pkg/api"
)

func createApiDevicesEndpoints(g *echo.Group, a *api.API) {
	g.GET("/devices", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, a.Devices)
	})
}
