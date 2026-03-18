package utils

import (
	"errors"
	"strconv"
	"strings"

	"github.com/knackwurstking/picow-led/pkg/models"
	"github.com/labstack/echo/v4"
)

var (
	ErrNotFound = errors.New("not found")
	ErrInvalid  = errors.New("invalid value")
)

func ParseQueryID(c echo.Context) (models.ID, error) {
	idStr := c.QueryParam("id")
	if idStr == "" {
		return 0, ErrNotFound
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, ErrInvalid
	}
	return models.ID(id), nil
}

func ParseQueryColor(c echo.Context) ([]uint8, error) {
	colorStr := c.QueryParam("color")
	if colorStr == "" {
		return nil, ErrNotFound
	}

	var color []uint8

	// Need to check value type first, e.g: "255,255,255" or "#FFFFFF"
	if strings.HasPrefix(colorStr, "#") {
		colorStr = strings.TrimPrefix(colorStr, "#")
		if len(colorStr) != 6 {
			return nil, ErrInvalid
		}
		for i := 0; i < 6; i += 2 {
			c, err := strconv.ParseUint(colorStr[i:i+2], 16, 8)
			if err != nil {
				return nil, ErrInvalid
			}
			color = append(color, uint8(c))
		}
	} else {
		colorStrSplit := strings.Split(colorStr, ",")
		for _, part := range colorStrSplit {
			part = strings.TrimSpace(part)
			if part == "" {
				color = append(color, 0)
				continue
			}
			c, err := strconv.ParseUint(part, 10, 8)
			if err != nil {
				return nil, ErrInvalid
			}
			color = append(color, uint8(c))
		}
	}

	return color, nil
}

func ParseQueryWhite(c echo.Context) (uint8, error) {
	whiteStr := c.QueryParam("white")
	if whiteStr == "" {
		return 0, ErrNotFound
	}

	white, err := strconv.ParseUint(whiteStr, 10, 8)
	if err != nil {
		return 0, ErrInvalid
	}

	return uint8(white), nil
}

func ParseQueryWhite2(c echo.Context) (uint8, error) {
	whiteStr := c.QueryParam("white2")
	if whiteStr == "" {
		return 0, ErrNotFound
	}

	white, err := strconv.ParseUint(whiteStr, 10, 8)
	if err != nil {
		return 0, ErrInvalid
	}

	return uint8(white), nil
}

func ParseParamID(c echo.Context) (models.ID, error) {
	idStr := c.Param("id")
	if idStr == "" {
		return 0, ErrNotFound
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, ErrInvalid
	}
	return models.ID(id), nil
}
