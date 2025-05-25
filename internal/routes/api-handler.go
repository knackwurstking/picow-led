package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"picow-led/internal/database"

	"github.com/labstack/echo/v4"
)

var (
	ErrorUnderConstruction = errors.New("under construction")
)

type ErrorResponse struct {
	Message error `json:"message"`
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Message: err,
	}
}

type APIHandler struct {
	db *database.DB
}

func NewAPIHandler(db *database.DB) *APIHandler {
	return &APIHandler{
		db: db,
	}
}

func (h *APIHandler) GetDevices(c echo.Context) error {
	// TODO: ...

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) GetDevicesAddr(c echo.Context) error {
	// TODO: ...

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) GetDevicesAddrName(c echo.Context) error {
	// TODO: ...

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) GetDevicesAddrActiveColor(c echo.Context) error {
	// TODO: ...

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) GetDevicesAddrColor(c echo.Context) error {
	// TODO: ...

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) GetDevicesAddrPins(c echo.Context) error {
	// TODO: ...

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) GetDevicesAddrPower(c echo.Context) error {
	// TODO: ...

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) PostDevicesAddrPower(c echo.Context) error {
	// TODO: ...

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) GetColors(c echo.Context) error {
	colors := h.db.Colors.List()
	data, err := json.Marshal(colors)
	if err != nil {
		return h.error(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, data)
}

func (h *APIHandler) PostColors(c echo.Context) error {
	// TODO: ...

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) PutColors(c echo.Context) error {
	// TODO: ...

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) GetColorsIndex(c echo.Context) error {
	// TODO: ...

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) PutColorsIndex(c echo.Context) error {
	// TODO: ...

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) DeleteColorsIndex(c echo.Context) error {
	// TODO: ...

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) error(c echo.Context, code int, err error) error {
	return c.JSON(code, NewErrorResponse(err))
}
