package utils

import (
	"fmt"
	"strconv"

	"github.com/knackwurstking/picow-led/pkg/models"
	"github.com/labstack/echo/v4"
)

func ParseQueryID(c echo.Context) (models.ID, error) {
	idStr := c.QueryParam("id")
	if idStr == "" {
		return 0, fmt.Errorf("missing device ID")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid device ID: %v", err)
	}
	return models.ID(id), nil
}

func ParseParamID(c echo.Context) (models.ID, error) {
	idStr := c.Param("id")
	if idStr == "" {
		return 0, fmt.Errorf("missing device ID")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid device ID: %v", err)
	}
	return models.ID(id), nil
}
