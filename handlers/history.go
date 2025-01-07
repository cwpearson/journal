package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/cwpearson/journal/database"
	"github.com/cwpearson/journal/models"
)

func HistoryGet(c echo.Context) error {
	db := database.Get()

	tagID := c.QueryParam("tag")

	var entries []models.Entry
	var err error
	var tag models.Tag

	if tagID == "" {
		// Retrieve all entries ordered by date
		err = db.Preload("Tags").
			Order("year DESC, month DESC, day DESC").
			Find(&entries).Error
	} else {
		// Retrieve entries ordered by date
		err = db.Joins("JOIN entry_tags ON entries.id = entry_tags.entry_id").
			Where("entry_tags.tag_id = ?", tagID).
			Order("entries.year DESC, entries.month DESC, entries.day DESC").
			Preload("Tags").
			Find(&entries).Error

		db.Where("id = ?", tagID).First(&tag)
	}

	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Database error")
	}
	return c.Render(http.StatusOK, "history.html", map[string]interface{}{
		"entries": entries,
		"tag":     tag.S,
	})
}
