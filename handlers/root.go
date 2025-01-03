package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RootGet(c echo.Context) error {
	return c.Render(http.StatusOK, "root.html", nil)
}
