package urlb

import (
	"fmt"
	"strings"

	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/pkg/models"
)

func APISetDeviceRGBW(deviceID models.ID, color []uint8, white uint8) string {
	colorStr := strings.Builder{}
	for i, c := range color {
		if i > 0 {
			colorStr.WriteString(",")
		}
		colorStr.WriteString(fmt.Sprintf("%d", c))
	}

	query := strings.Builder{}
	query.WriteString("?color=")
	query.WriteString(colorStr.String())

	query.WriteString(fmt.Sprintf("&white=%d", white))

	return env.Route(fmt.Sprintf("/api/device/%d/color%s", deviceID, query.String()))
}

func APISetDeviceWhite(deviceID models.ID, white uint8) string {
	query := fmt.Sprintf("?white=%d", white)
	return env.Route(fmt.Sprintf("/api/device/%d/white%s", deviceID, query))
}
