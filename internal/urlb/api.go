package urlb

import (
	"fmt"
	"strings"

	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/pkg/models"
)

func APISetDeviceColor(deviceID models.ID, color ...uint8) string {
	colorStr := strings.Builder{}
	for i, c := range color {
		if i > 0 {
			colorStr.WriteString(",")
		}
		colorStr.WriteString(fmt.Sprintf("%d", c))
	}

	query := strings.Builder{}
	if colorStr.Len() > 0 {
		query.WriteString("?color=")
		query.WriteString(colorStr.String())
	}

	return env.Route(fmt.Sprintf("/api/device/%d/color%s", deviceID, query.String()))
}
