package assets

import (
	"embed"

	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/labstack/echo/v4"
)

var (
	//go:embed public
	public embed.FS
)

func ServePublicFS(e *echo.Echo) *echo.Route {
	fs := echo.MustSubFS(public, "public")
	return e.StaticFS(env.Route("/"), fs)
}
