package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/cwpearson/journal/database"
	"github.com/cwpearson/journal/models"
)

func TagGet(c echo.Context) error {
	db := database.Get()

	tagID := c.QueryParam("tag")
	if tagID == "" {
		return c.String(http.StatusInternalServerError, "no tag provided")
	}

	// Retrieve entries ordered by date
	var entries []models.Entry
	err := db.Joins("JOIN entry_tags ON entries.id = entry_tags.entry_id").
		Where("entry_tags.tag_id = ?", tagID).
		Order("entries.year DESC, entries.month DESC, entries.day DESC").
		Preload("Tags").
		Find(&entries).Error
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "database error")
	}

	// Return HTML template
	return c.Render(http.StatusOK, "history.html", entries)
}
