package urlb

import (
	"fmt"
	"strings"

	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/pkg/models"
)

func APISetDeviceColor(deviceID models.ID, color ...uint8) string {
	colorStr := &strings.Builder{}
	for i, c := range color {
		if i > 0 {
			colorStr.WriteString(",")
		}
		fmt.Fprintf(colorStr, "%d", c)
	}

	query := &strings.Builder{}
	if len(color) > 0 {
		query.WriteString("?color=")
		query.WriteString(colorStr.String())
	}

	return env.Route(fmt.Sprintf("/api/device/%d/color%s", deviceID, query.String()))
}

func APISetDeviceWhite(deviceID models.ID, white uint8) string {
	query := fmt.Sprintf("?white=%d", white)
	return env.Route(fmt.Sprintf("/api/device/%d/white%s", deviceID, query))
}

func APISetDeviceRGBW(deviceID models.ID, color []uint8, white uint8) string {
	colorStr := &strings.Builder{}
	for i, c := range color {
		if i > 0 {
			colorStr.WriteString(",")
		}
		fmt.Fprintf(colorStr, "%d", c)
	}

	query := &strings.Builder{}
	query.WriteString("?color=")
	query.WriteString(colorStr.String())

	fmt.Fprintf(query, "&white=%d", white)

	return env.Route(fmt.Sprintf("/api/device/%d/color%s", deviceID, query.String()))
}
