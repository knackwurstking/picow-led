package urlb

import (
	"fmt"

	"github.com/knackwurstking/picow-led/internal/env"
	"github.com/knackwurstking/picow-led/pkg/models"
)

func PageHome() string {
	return env.Route("/")
}

func PageDevice(deviceID models.ID) string {
	return env.Route(fmt.Sprintf("/device?id=%d", deviceID))
}
