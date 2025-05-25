package routes

import (
	"errors"

	"github.com/labstack/echo/v4"
)

var ErrorUnderConstruction = errors.New("under construction")

type APIHandler struct{}

func NewAPIHandler() *APIHandler {
	return &APIHandler{}
}

func (h *APIHandler) GetDevices(c echo.Context) error {
	// TODO: ...

	return ErrorUnderConstruction
}

func (h *APIHandler) GetDevicesAddr(c echo.Context) error {
	// TODO: ...

	return ErrorUnderConstruction
}

func (h *APIHandler) GetDevicesAddrName(c echo.Context) error {
	// TODO: ...

	return ErrorUnderConstruction
}

func (h *APIHandler) GetDevicesAddrActiveColor(c echo.Context) error {
	// TODO: ...

	return ErrorUnderConstruction
}

func (h *APIHandler) GetDevicesAddrColor(c echo.Context) error {
	// TODO: ...

	return ErrorUnderConstruction
}

func (h *APIHandler) GetDevicesAddrPins(c echo.Context) error {
	// TODO: ...

	return ErrorUnderConstruction
}

func (h *APIHandler) GetDevicesAddrPower(c echo.Context) error {
	// TODO: ...

	return ErrorUnderConstruction
}

func (h *APIHandler) PostDevicesAddrPower(c echo.Context) error {
	// TODO: ...

	return ErrorUnderConstruction
}

func (h *APIHandler) GetColors(c echo.Context) error {
	// TODO: ...

	return ErrorUnderConstruction
}

func (h *APIHandler) PostColors(c echo.Context) error {
	// TODO: ...

	return ErrorUnderConstruction
}

func (h *APIHandler) PutColors(c echo.Context) error {
	// TODO: ...

	return ErrorUnderConstruction
}

func (h *APIHandler) GetColorsIndex(c echo.Context) error {
	// TODO: ...

	return ErrorUnderConstruction
}

func (h *APIHandler) PutColorsIndex(c echo.Context) error {
	// TODO: ...

	return ErrorUnderConstruction
}

func (h *APIHandler) DeleteColorsIndex(c echo.Context) error {
	// TODO: ...

	return ErrorUnderConstruction
}
