package endpoints

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

func readBodyData(c echo.Context, data any) (status int, err error) {
	d, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if err := json.Unmarshal(d, data); err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func saveConfig(config *Config, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := config.save(); err != nil {
			slog.Warn("Save config", "path", config.Path, "error", err)
		}
	}()
	go config.save()
}
