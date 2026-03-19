package handlers

import (
	"fmt"
	"net/http"

	"github.com/knackwurstking/picow-led/internal/services"
	"github.com/labstack/echo/v4"
)

func handleServiceError(err error, message string) *echo.HTTPError {
	status := http.StatusInternalServerError
	switch err {
	case services.ErrorNotFound:
		status = http.StatusNotFound
	case services.ErrorValidation:
		status = http.StatusBadRequest
	}
	return echo.NewHTTPError(status, fmt.Sprintf("%s: %s", message, err.Error()))
}
