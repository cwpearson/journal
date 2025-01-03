package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/cwpearson/journal/database"
	"github.com/cwpearson/journal/models"
)

func HistoryGet(c echo.Context) error {
	db := database.Get()

	var entries []models.Entry

	// Retrieve all entries ordered by date
	result := db.Preload("Tags").
		Order("year DESC, month DESC, day DESC").
		Find(&entries)

	if result.Error != nil {
		return c.String(http.StatusInternalServerError, "Database error")
	}

	// Return HTML template
	return c.Render(http.StatusOK, "history.html", entries)
}
