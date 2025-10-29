package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/knackwurstking/picow-led/env"
	"github.com/knackwurstking/picow-led/models"
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, method, path string, handler echo.HandlerFunc) {
	switch method {
	case http.MethodGet:
		e.GET(env.Args.ServerPathPrefix+path, handler)
		e.GET(env.Args.ServerPathPrefix+path+"/", handler)
	case http.MethodPost:
		e.POST(env.Args.ServerPathPrefix+path, handler)
		e.POST(env.Args.ServerPathPrefix+path+"/", handler)
	case http.MethodPut:
		e.PUT(env.Args.ServerPathPrefix+path, handler)
		e.PUT(env.Args.ServerPathPrefix+path+"/", handler)
	case http.MethodDelete:
		e.DELETE(env.Args.ServerPathPrefix+path, handler)
		e.DELETE(env.Args.ServerPathPrefix+path+"/", handler)
	default:
		panic(fmt.Sprintf("unsupported method: %s", method))
	}
}

func QueryParam(c echo.Context, paramName string, optional bool) (string, error) {
	if paramQuery := c.QueryParam(paramName); paramQuery != "" {
		return paramQuery, nil
	}

	if optional {
		return "", nil
	}

	return "", NewValidationError("missing %s query parameter", paramName)
}

func QueryParamDeviceID(c echo.Context, paramName string, optional bool) (models.DeviceID, error) {
	param, err := QueryParam(c, paramName, optional)
	if err != nil {
		return 0, err
	}

	if optional && param == "" {
		return 0, nil
	}

	deviceIDConversion, err := strconv.Atoi(param)
	if err != nil {
		return 0, NewValidationError("invalid device ID query parameter: %s=%s", paramName, param)
	}
	return models.DeviceID(deviceIDConversion), nil
}
