package endpoints

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/knackwurstking/picow-led-server/pkg/api"
)

func createApiColorsEndpoints(g *echo.Group, a *api.API, changeCallback func()) {
	g.GET("/colors", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, a.Colors)
	})

	g.GET("/colors/:name", func(ctx echo.Context) error {
		name := ctx.Param("name")

		for k, v := range a.Colors {
			if k == name {
				return ctx.JSON(http.StatusOK, v)
			}
		}

		return ctx.String(http.StatusBadRequest, fmt.Sprintf("color \"%s\" not found", name))
	})

	g.POST("/colors/:name", func(ctx echo.Context) error {
		name := ctx.Param("name")
		reqData := make(api.Color, 0)
		if status, err := readBodyData(ctx, &reqData); err != nil {
			return ctx.String(status, err.Error())
		}

		if err := a.AddColor(name, reqData); err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}

		defer changeCallback()
		return ctx.JSON(http.StatusOK, nil)
	})

	g.PUT("/colors/:name", func(ctx echo.Context) error {
		name := ctx.Param("name")

		reqData := make(api.Color, 0)
		if status, err := readBodyData(ctx, &reqData); err != nil {
			return ctx.String(status, err.Error())
		}

		if err := a.ReplaceColor(name, reqData); err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}

		defer changeCallback()
		return ctx.JSON(http.StatusOK, nil)
	})

	g.DELETE("/colors/:name", func(ctx echo.Context) error {
		name := ctx.Param("name")
		a.DeleteColor(name)
		defer changeCallback()
		return ctx.JSON(http.StatusOK, nil)
	})
}
