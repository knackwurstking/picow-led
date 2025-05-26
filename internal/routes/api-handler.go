package routes

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"picow-led/internal/database"
	"strconv"

	"github.com/labstack/echo/v4"
)

var ErrorUnderConstruction = errors.New("under construction")

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{
		Message: msg,
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
	devices, err := h.db.Devices.List()
	if err != nil {
		return h.error(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, devices)
}

func (h *APIHandler) GetDevice(c echo.Context) error {
	addr, err := url.QueryUnescape(c.Param("addr"))
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	device, err := h.db.Devices.Get(addr)
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, device)
}

func (h *APIHandler) GetDeviceName(c echo.Context) error {
	// TODO: Continue here

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) GetDeviceActiveColor(c echo.Context) error {
	// TODO: ...

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) GetDeviceColor(c echo.Context) error {
	// TODO: ...

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) GetDevicePins(c echo.Context) error {
	// TODO: ...

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) GetDevicePower(c echo.Context) error {
	// TODO: ...

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) PostDevicePower(c echo.Context) error {
	// TODO: ...

	return h.error(c, http.StatusInternalServerError, ErrorUnderConstruction)
}

func (h *APIHandler) GetColors(c echo.Context) error {
	colors, err := h.db.Colors.List()
	if err != nil {
		return h.error(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, colors)
}

func (h *APIHandler) PostColors(c echo.Context) error {
	colors := []database.Color{}
	err := json.NewDecoder(c.Request().Body).Decode(&colors)
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	err = h.db.Colors.Add(colors...)
	if err != nil {
		return h.error(c, http.StatusInternalServerError, err)
	}

	return nil
}

func (h *APIHandler) PutColors(c echo.Context) error {
	colors := []database.Color{}
	err := json.NewDecoder(c.Request().Body).Decode(&colors)
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	err = h.db.Colors.Set(colors...)
	if err != nil {
		return h.error(c, http.StatusInternalServerError, err)
	}

	return nil
}

func (h *APIHandler) GetColorsID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	color, err := h.db.Colors.Get(id)
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, color)
}

func (h *APIHandler) PostColorsID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	color := database.Color{}
	err = json.NewDecoder(c.Request().Body).Decode(&color)
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	err = h.db.Colors.Update(id, color)
	if err != nil {
		return h.error(c, http.StatusInternalServerError, err)
	}

	return err
}

func (h *APIHandler) DeleteColorsID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	err = h.db.Colors.Delete(id)
	if err != nil {
		return h.error(c, http.StatusInternalServerError, err)
	}

	return nil
}

func (h *APIHandler) error(c echo.Context, code int, err error) error {
	slog.Error(err.Error())
	return c.JSON(code, NewErrorResponse(err.Error()))
}
