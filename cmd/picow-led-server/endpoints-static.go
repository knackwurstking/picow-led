package main

import (
	"net/http"

	"github.com/knackwurstking/picow-led-server/frontend"
	"github.com/labstack/echo/v4"
)

func endpointsStatic() {
	fs := frontend.GetFS()
	e.StaticFS("/", fs)
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusSeeOther, "/index.html")
	})
}
