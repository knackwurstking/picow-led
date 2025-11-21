package dialogs

import (
	"net/http"

	"github.com/knackwurstking/picow-led/services"
	"github.com/knackwurstking/picow-led/utils"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	registry *services.Registry
}

func NewHandler(r *services.Registry) *Handler {
	return &Handler{
		registry: r,
	}
}

func (h *Handler) Register(e *echo.Echo) {
	// Edit Device
	utils.Register(e, http.MethodGet, "/htmx/dialog/edit-device", h.GetEditDevice)
	utils.Register(e, http.MethodPost, "/htmx/dialog/edit-device", h.PostEditDevice)
	utils.Register(e, http.MethodPut, "/htmx/dialog/edit-device", h.PutEditDevice)

	// Edit Group
	utils.Register(e, http.MethodGet, "/htmx/dialog/edit-group", h.GetEditGroup)
	utils.Register(e, http.MethodPost, "/htmx/dialog/edit-group", h.PostEditGroup)
	utils.Register(e, http.MethodPut, "/htmx/dialog/edit-group", h.PutEditGroup)
}
