package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func endpointsApi() {
	g := e.Group("/api")

	g.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, api)
	})

	setupEndpointApiDevices(g)
	setupEndpointApiColors(g)
}
