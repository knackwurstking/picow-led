package routes

import (
	"picow-led/components"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func frontend(e *echo.Echo, data Options) {
	// Base Data (templ)
	baseData := &components.BaseData{
		ServerPathPrefix: data.ServerPathPrefix,
		Version:          data.Version,
	}

	// Page Data: "/" - devices
	pageDevicesData := &components.PageDevicesData{
		BaseData: baseData,
		Devices:  data.Api.Servers, // TODO: To fix this is need an global api storage first
	}

	e.GET(data.ServerPathPrefix+"/", func(c echo.Context) error {
		return renderTempl(c,
			components.Base(baseData,
				components.PageDevices(pageDevicesData),
			),
		)
	})

	e.GET(data.ServerPathPrefix+"/settings", func(c echo.Context) error {
		return renderTempl(c,
			components.Base(baseData,
				components.PageSettings(),
			),
		)
	})
}

func renderTempl(c echo.Context, t templ.Component) error {
	return t.Render(c.Request().Context(), c.Response().Writer)
}
