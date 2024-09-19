package endpoints

import (
	"github.com/labstack/echo/v4"

	_clients "github.com/knackwurstking/picow-led-server/pkg/clients"
)

// TODO: Add some "(api)change" callback (param) to Create function
//func saveConfig() {
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		if err := config.save(); err != nil {
//			slog.Warn("Save config", "config.Path", config.Path, "err", err)
//		}
//	}()
//	go config.save()
//}

// TODO: Missing *api.API, need to passed to `createApiEndpoints`
func Create(e *echo.Echo) {
	clients := _clients.NewClients()

	createStaticEndpoints(e)
	createEventsEndpoints(e, clients)
	createApiEndpoints(e)
}
