package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/cwpearson/journal/config"
	"github.com/cwpearson/journal/database"
	"github.com/cwpearson/journal/entries"
	"github.com/cwpearson/journal/models"
	"github.com/cwpearson/journal/ollama"
	"github.com/labstack/echo/v4"
)

func EditGet(c echo.Context) error {

	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}
	id := uint(id64)

	// retrieve text for entry
	db := database.Get()

	var entry models.Entry
	err = db.Where("id = ?", id).Find(&entry).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}

	year, month, day, name := entries.PathComponents(&entry)
	textPath := filepath.Join(config.DataDir(), year, month, day, name)

	log.Println("read", textPath)
	contents, err := os.ReadFile(textPath)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}

	return c.Render(http.StatusOK, "edit.html", map[string]interface{}{
		"title":   "Edit Entry",
		"mode":    "edit",
		"text":    string(contents),
		"entryId": id,
	})

}

func EditPost(c echo.Context) error {
	log.Println("EditPost")

	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}
	id := uint(id64)

	submission := new(CreateForm)

	// Bind form data to struct
	if err := c.Bind(submission); err != nil {
		return c.String(http.StatusBadRequest, "Invalid submission")
	}

	// Validate submission
	if submission.UserText == "" {
		return c.String(http.StatusBadRequest, "Text cannot be empty")
	}

	if err = entries.Set(id, submission.UserText, ollama.NewClientFromConfig()); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}

	referer := c.Request().Referer()
	if referer == "" {
		referer = "/history"
	}
	return c.Redirect(http.StatusSeeOther, referer)
}
