package handlers

import (
	"fmt"
	"net/http"

	"github.com/knackwurstking/picow-led/env"
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
