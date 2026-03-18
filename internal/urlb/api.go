package urlb

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/pkg/models"
)

func APISetDeviceColor(deviceID models.ID, color ...uint8) string {
	v := url.Values{}

	if len(color) > 0 {
		colorStr := &strings.Builder{}
		for i, c := range color {
			if i > 0 {
				colorStr.WriteString(",")
			}
			fmt.Fprintf(colorStr, "%d", c)
		}
		v.Add("color", colorStr.String())
	}

	return env.Route(fmt.Sprintf("/api/devices/%d/color%s", deviceID, v.Encode()))
}

func APISetDeviceWhite(deviceID models.ID, white uint8) string {
	v := url.Values{}
	v.Add("white", fmt.Sprintf("%d", white))
	return env.Route(fmt.Sprintf("/api/devices/%d/white%s", deviceID, v.Encode()))
}

func APISetDeviceRGBW(deviceID models.ID, color []uint8, white uint8) string {
	v := url.Values{}

	if len(color) > 0 {
		colorStr := &strings.Builder{}
		for i, c := range color {
			if i > 0 {
				colorStr.WriteString(",")
			}
			fmt.Fprintf(colorStr, "%d", c)
		}
		v.Add("color", colorStr.String())
	}

	v.Add("white", fmt.Sprintf("%d", white))

	return env.Route(fmt.Sprintf("/api/devices/%d/color%s", deviceID, v.Encode()))
}
