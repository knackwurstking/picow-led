package endpoints

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"

	"github.com/knackwurstking/picow-led-server/pkg/api"
	"github.com/knackwurstking/picow-led-server/pkg/clients"
)

func Create(e *echo.Echo) {
	c := clients.NewClients()
	a := api.NewAPI()

	config := NewConfig(a)
	slog.Debug(fmt.Sprintf("Try to load API configuration from %s", config.Path))
	if err := config.load(); err != nil {
		panic(fmt.Sprintf("Loading API configuration froom \"%s\": %s", config.Path, err))
	}

	createStaticEndpoints(e)
	createEventsEndpoints(e, c)

	createApiEndpoints(e, c, a, func() {
		if err := config.save(); err != nil {
			slog.Warn(fmt.Sprintf("Save API configuration: %s", err))
		}
	})
}
