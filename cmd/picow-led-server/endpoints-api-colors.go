package main

import (
	"fmt"
	"net/http"

	_api "github.com/knackwurstking/picow-led-server/pkg/api"
	"github.com/labstack/echo/v4"
)

func setupEndpointApiColors(g *echo.Group) {
	g.GET("/colors", func(c echo.Context) error {
		return c.JSON(http.StatusOK, api.Colors)
	})

	g.GET("/colors/:name", func(c echo.Context) error {
		name := c.Param("name")

		for k, v := range api.Colors {
			if k == name {
				return c.JSON(http.StatusOK, v)
			}
		}

		return c.JSON(http.StatusBadRequest, fmt.Sprintf("color \"%s\" not found", name))
	})

	g.POST("/colors/:name", func(c echo.Context) error {
		name := c.Param("name")

		r := make(_api.Color, 0)
		status, err := readBodyData(c, &r)
		if status != http.StatusOK {
			return c.JSON(status, err.Error())
		}

		err = api.Colors.Add(name, r)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		saveConfig()
		return c.JSON(http.StatusOK, nil)
	})

	g.PUT("/colors/:name", func(c echo.Context) error {
		name := c.Param("name")

		r := make(_api.Color, 0)
		status, err := readBodyData(c, &r)
		if status != http.StatusOK {
			return c.JSON(status, err.Error())
		}

		err = api.Colors.Replace(name, r)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		saveConfig()
		return c.JSON(http.StatusOK, nil)
	})

	g.DELETE("/colors/:name", func(c echo.Context) error {
		name := c.Param("name")
		api.Colors.Remove(name)
		saveConfig()
		return c.JSON(http.StatusOK, nil)
	})
}
