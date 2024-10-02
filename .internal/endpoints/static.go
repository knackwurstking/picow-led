package endpoints

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/knackwurstking/picow-led-server/frontend"
)

func createStaticEndpoints(e *echo.Echo) {
	fs := frontend.GetFS()
	e.StaticFS("/", fs)
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusSeeOther, "/index.html")
	})
}
