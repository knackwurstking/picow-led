package routes

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"picow-led/internal/database"
	"picow-led/internal/micro"
	"picow-led/internal/ws"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{
		Message: msg,
	}
}

type APIHandler struct {
	db        *database.DB
	wsHandler *ws.Handler
}

func NewAPIHandler(db *database.DB, wsHandler *ws.Handler) *APIHandler {
	return &APIHandler{
		db:        db,
		wsHandler: wsHandler,
	}
}

func (h *APIHandler) GetDevices(c echo.Context) error {
	devices, err := h.db.Devices.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, devices)
}

func (h *APIHandler) GetDevice(c echo.Context) error {
	addr, err := url.QueryUnescape(c.Param("addr"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	device, err := h.db.Devices.Get(addr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, device)
}

func (h *APIHandler) GetDeviceName(c echo.Context) error {
	addr, err := url.QueryUnescape(c.Param("addr"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	device, err := h.db.Devices.Get(addr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, device.Name)
}

func (h *APIHandler) GetDeviceActiveColor(c echo.Context) error {
	addr, err := url.QueryUnescape(c.Param("addr"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	device, err := h.db.Devices.Get(addr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, device.ActiveColor)
}

func (h *APIHandler) GetDeviceColor(c echo.Context) error {
	addr, err := url.QueryUnescape(c.Param("addr"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	device, err := h.db.Devices.Get(addr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, device.Color)
}

func (h *APIHandler) PostDeviceColor(c echo.Context) error {
	addr, err := url.QueryUnescape(c.Param("addr"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	device, err := h.db.Devices.Get(addr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	color := []int{}
	err = json.NewDecoder(c.Request().Body).Decode(&color)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = micro.SetColor(device.Addr, color)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	device.SetColor(color)
	err = h.db.Devices.Update(device.Addr, device)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	go h.wsHandler.BroadcastDevice(device)

	return nil
}

func (h *APIHandler) GetDevicePins(c echo.Context) error {
	addr, err := url.QueryUnescape(c.Param("addr"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	device, err := h.db.Devices.Get(addr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, device.Pins)
}

func (h *APIHandler) GetDevicePower(c echo.Context) error {
	addr, err := url.QueryUnescape(c.Param("addr"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	device, err := h.db.Devices.Get(addr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, device.Power)
}

func (h *APIHandler) PostDevicePower(c echo.Context) error {
	addr, err := url.QueryUnescape(c.Param("addr"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	state, err := strconv.Atoi(c.QueryParam("state"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	device, err := h.db.Devices.Get(addr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	color := []int{}
	switch database.PowerState(state) {
	case database.PowerStateOFF:
		color = device.GetColorForPowerState(database.PowerStateOFF)
		if err := micro.SetColor(device.Addr, color); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	case database.PowerStateON:
		color = device.GetColorForPowerState(database.PowerStateON)
		if err := micro.SetColor(device.Addr, color); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	default:
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf(
			"unknown power state %d, expect 1 (ON) or 0 (OFF)", state,
		))
	}

	device.SetColor(color)
	if err = h.db.Devices.Update(device.Addr, device); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	go h.wsHandler.BroadcastDevice(device)

	return nil
}

func (h *APIHandler) GetColors(c echo.Context) error {
	colors, err := h.db.Colors.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, colors)
}

func (h *APIHandler) PostColors(c echo.Context) error {
	colors := []database.Color{}
	err := json.NewDecoder(c.Request().Body).Decode(&colors)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = h.db.Colors.Add(colors...)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	go func() {
		if colors, err = h.db.Colors.List(); err != nil {
			slog.Error("Handle broadcast colors", "error", err)
		} else {
			h.wsHandler.BroadcastColors(colors)
		}
	}()

	return nil
}

func (h *APIHandler) PutColors(c echo.Context) error {
	colors := []database.Color{}
	err := json.NewDecoder(c.Request().Body).Decode(&colors)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = h.db.Colors.Set(colors...)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	go h.wsHandler.BroadcastColors(colors)

	return nil
}

func (h *APIHandler) GetColorsID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	color, err := h.db.Colors.Get(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, color)
}

func (h *APIHandler) PostColorsID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	color := database.Color{}
	err = json.NewDecoder(c.Request().Body).Decode(&color)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = h.db.Colors.Update(id, color)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	go func() {
		if colors, err := h.db.Colors.List(); err != nil {
			slog.Error("broadcast colors", "error", err)
		} else {
			h.wsHandler.BroadcastColors(colors)
		}
	}()

	return err
}

func (h *APIHandler) DeleteColorsID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = h.db.Colors.Delete(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	go func() {
		if colors, err := h.db.Colors.List(); err != nil {
			slog.Error("broadcast colors", "error", err)
		} else {
			h.wsHandler.BroadcastColors(colors)
		}
	}()

	return nil
}
