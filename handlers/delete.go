package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cwpearson/journal/entries"
	"github.com/labstack/echo/v4"
)

func DeletePost(c echo.Context) error {
	log.Println("DeletePost")

	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}
	id := uint(id64)

	if err = entries.Delete(id); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}

	return c.Redirect(http.StatusSeeOther, "/history")
}
