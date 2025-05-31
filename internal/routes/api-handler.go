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
		return h.error(c, http.StatusInternalServerError,
			fmt.Errorf("database: list devices: %s", err))
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
		return h.error(c, http.StatusBadRequest,
			fmt.Errorf("database: get device: %s", err))
	}

	return c.JSON(http.StatusOK, device)
}

func (h *APIHandler) GetDeviceName(c echo.Context) error {
	addr, err := url.QueryUnescape(c.Param("addr"))
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	device, err := h.db.Devices.Get(addr)
	if err != nil {
		return h.error(c, http.StatusBadRequest,
			fmt.Errorf("database: get device: %s", err))
	}

	return c.JSON(http.StatusOK, device.Name)
}

func (h *APIHandler) GetDeviceActiveColor(c echo.Context) error {
	addr, err := url.QueryUnescape(c.Param("addr"))
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	device, err := h.db.Devices.Get(addr)
	if err != nil {
		return h.error(c, http.StatusBadRequest,
			fmt.Errorf("database: get device: %s", err))
	}

	return c.JSON(http.StatusOK, device.ActiveColor)
}

func (h *APIHandler) GetDeviceColor(c echo.Context) error {
	addr, err := url.QueryUnescape(c.Param("addr"))
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	device, err := h.db.Devices.Get(addr)
	if err != nil {
		return h.error(c, http.StatusBadRequest,
			fmt.Errorf("database: get device: %s", err))
	}

	return c.JSON(http.StatusOK, device.Color)
}

// TODO: Handle ws event "device"
func (h *APIHandler) PostDeviceColor(c echo.Context) error {
	addr, err := url.QueryUnescape(c.Param("addr"))
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	device, err := h.db.Devices.Get(addr)
	if err != nil {
		return h.error(c, http.StatusBadRequest,
			fmt.Errorf("database: get device: %s", err))
	}

	color := []int{}
	err = json.NewDecoder(c.Request().Body).Decode(&color)
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	device.SetColor(color)
	err = micro.SetColor(device.Addr, device.Color)
	if err != nil {
		return h.error(c, http.StatusInternalServerError,
			fmt.Errorf("micro: set device color: %s", err))
	}

	err = h.db.Devices.Update(device.Addr, device)
	if err != nil {
		return h.error(c, http.StatusInternalServerError,
			fmt.Errorf("database: update device: %s", err))
	}

	return nil
}

func (h *APIHandler) GetDevicePins(c echo.Context) error {
	addr, err := url.QueryUnescape(c.Param("addr"))
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	device, err := h.db.Devices.Get(addr)
	if err != nil {
		return h.error(c, http.StatusBadRequest,
			fmt.Errorf("database: get device: %s", err))
	}

	return c.JSON(http.StatusOK, device.Pins)
}

func (h *APIHandler) GetDevicePower(c echo.Context) error {
	addr, err := url.QueryUnescape(c.Param("addr"))
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	device, err := h.db.Devices.Get(addr)
	if err != nil {
		return h.error(c, http.StatusBadRequest,
			fmt.Errorf("database: get device: %s", err))
	}

	return c.JSON(http.StatusOK, device.Power)
}

// TODO: Handle ws event "device"
func (h *APIHandler) PostDevicePower(c echo.Context) error {
	addr, err := url.QueryUnescape(c.Param("addr"))
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	state, err := strconv.Atoi(c.QueryParam("state"))
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	device, err := h.db.Devices.Get(addr)
	if err != nil {
		return h.error(c, http.StatusBadRequest,
			fmt.Errorf("database: get device: %s", err))
	}

	color := []int{}
	switch database.PowerState(state) {
	case database.PowerStateOFF:
		color = device.GetColorForPowerState(database.PowerStateOFF)
		if err := micro.SetColor(device.Addr, color); err != nil {
			return h.error(c, http.StatusInternalServerError,
				fmt.Errorf("micro: set color: %s", err))
		}
	case database.PowerStateON:
		color = device.GetColorForPowerState(database.PowerStateON)
		if err := micro.SetColor(device.Addr, color); err != nil {
			return h.error(c, http.StatusInternalServerError,
				fmt.Errorf("micro: set color: %s", err))
		} else {
			device.SetColor(color)
		}
	default:
		return h.error(c, http.StatusBadRequest,
			fmt.Errorf(
				"unknown power state %d, expect 1 (ON) or 0 (OFF)", state,
			),
		)
	}

	device.SetColor(color)
	if err = h.db.Devices.Update(device.Addr, device); err != nil {
		return h.error(c, http.StatusInternalServerError,
			fmt.Errorf("database: update device: %s", err))
	}

	return nil
}

func (h *APIHandler) GetColors(c echo.Context) error {
	colors, err := h.db.Colors.List()
	if err != nil {
		return h.error(c, http.StatusInternalServerError,
			fmt.Errorf("database: get colors: %s", err))
	}
	return c.JSON(http.StatusOK, colors)
}

// TODO: Handle ws event "colors"
func (h *APIHandler) PostColors(c echo.Context) error {
	colors := []database.Color{}
	err := json.NewDecoder(c.Request().Body).Decode(&colors)
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	err = h.db.Colors.Add(colors...)
	if err != nil {
		return h.error(c, http.StatusInternalServerError,
			fmt.Errorf("database: add colors: %s", err))
	}

	return nil
}

// TODO: Handle ws event "colors"
func (h *APIHandler) PutColors(c echo.Context) error {
	colors := []database.Color{}
	err := json.NewDecoder(c.Request().Body).Decode(&colors)
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	err = h.db.Colors.Set(colors...)
	if err != nil {
		return h.error(c, http.StatusInternalServerError,
			fmt.Errorf("database: set colors: %s", err))
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
		return h.error(c, http.StatusBadRequest,
			fmt.Errorf("database: get color: %s", err))
	}

	return c.JSON(http.StatusOK, color)
}

// TODO: Handle ws event "color"
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
		return h.error(c, http.StatusInternalServerError,
			fmt.Errorf("database: update color: %s", err))
	}

	return err
}

// TODO: Handle ws event "color"
func (h *APIHandler) DeleteColorsID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	err = h.db.Colors.Delete(id)
	if err != nil {
		return h.error(c, http.StatusInternalServerError,
			fmt.Errorf("database: delete color: %s", err))
	}

	return nil
}

func (h *APIHandler) error(c echo.Context, code int, err error) error {
	slog.Error(err.Error())
	return c.JSON(code, NewErrorResponse(err.Error()))
}
