//go:build frontend
// +build frontend

package main

import (
	"embed"

	"github.com/labstack/echo/v4"
)

//go:embed frontend-build
var frontendBuild embed.FS

func StaticFS(e *echo.Echo) {
	e.StaticFS(serverPathPrefix+"/", echo.MustSubFS(frontendBuild, "frontend-build"))
}
