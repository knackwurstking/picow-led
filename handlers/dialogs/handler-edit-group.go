package dialogs

import "github.com/labstack/echo/v4"

func (h *Handler) GetEditGroup(c echo.Context) error {
	// TODO: ...

	return nil
}

func (h *Handler) PostEditGroup(c echo.Context) error {
	// TODO: ...

	c.Response().Header().Set("HX-Trigger", "reloadGroups")
	return nil
}
func (h *Handler) PutEditGroup(c echo.Context) error {
	// TODO: ...

	c.Response().Header().Set("HX-Trigger", "reloadGroups")
	return nil
}
