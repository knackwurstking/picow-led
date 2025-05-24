package routes

import "picow-led/internal/types"

type PageControl struct {
	Global
	Device *types.Device
}
