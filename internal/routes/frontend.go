package routes

import (
	"picow-led/components"
	"picow-led/internal/api"
	"picow-led/ui"

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
	devices := []*api.Device{}
	for _, s := range data.Api.Servers {
		devices = append(devices, &api.Device{Server: s})
	}
	pageDevicesData := &components.PageDevicesData{
		BaseData: baseData,
		Devices:  devices,
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

	e.GET(data.ServerPathPrefix+"/gomponents", func(c echo.Context) error {
		return ui.DevicesPage().Render(c.Response().Writer)
	})
}

func renderTempl(c echo.Context, t templ.Component) error {
	return t.Render(c.Request().Context(), c.Response().Writer)
}
