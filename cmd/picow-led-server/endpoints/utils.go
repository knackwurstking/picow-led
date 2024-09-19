package endpoints

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func readBodyData(c echo.Context, data any) (status int, err error) {
	d, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = json.Unmarshal(d, data)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}
