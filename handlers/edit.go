package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/cwpearson/journal/database"
	"github.com/cwpearson/journal/entries"
	"github.com/cwpearson/journal/models"
	"github.com/cwpearson/journal/ollama"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func parseEditParams(c echo.Context) (int, int, int, error) {
	yStr, mStr, dStr := c.Param("year"), c.Param("month"), c.Param("day")

	y64, err := strconv.ParseInt(yStr, 10, 32)
	if err != nil {
		return 0, 0, 0, err
	}
	m64, err := strconv.ParseInt(mStr, 10, 32)
	if err != nil {
		return 0, 0, 0, err
	}
	d64, err := strconv.ParseInt(dStr, 10, 32)
	if err != nil {
		return 0, 0, 0, err
	}

	return int(y64), int(m64), int(d64), nil
}

func formatDate(year int, month int, day int) string {
	// Convert month number to abbreviated month name
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	monthStr := months[month-1]

	// Get ordinal suffix for day
	suffix := "th"
	if day%10 == 1 && day != 11 {
		suffix = "st"
	} else if day%10 == 2 && day != 12 {
		suffix = "nd"
	} else if day%10 == 3 && day != 13 {
		suffix = "rd"
	}

	// Format the full date string
	return fmt.Sprintf("%s %d%s, %d", monthStr, day, suffix, year)
}

func EditGet(c echo.Context) error {
	year, month, day, err := parseEditParams(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}
	log.Println("EditGet", year, month, day)

	// retrieve text for entry
	db := database.Get()

	entry := new(models.Entry)
	err = db.
		Where("year = ?", year).
		Where("month = ?", month).
		Where("day = ?", day).
		First(entry).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			entry, err = entries.Create(year, month, day, "")
			if err != nil {
				return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
			}
		} else {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
		}
	}

	textPath := entries.TextPath(entry)
	log.Println("read", textPath)
	contents, err := os.ReadFile(textPath)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}

	return c.Render(http.StatusOK, "edit.html", map[string]interface{}{
		"title":   formatDate(year, month, day),
		"year":    year,
		"month":   month,
		"day":     day,
		"text":    string(contents),
		"entryId": entry.ID,
	})

}

type EditForm struct {
	UserText string `form:"userText"`
}

func EditPost(c echo.Context) error {
	log.Println("EditPost")

	year, month, day, err := parseEditParams(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}
	log.Println("EditGet", year, month, day)

	// retrieve text for entry
	db := database.Get()

	entry := new(models.Entry)
	err = db.
		Where("year = ?", year).
		Where("month = ?", month).
		Where("day = ?", day).
		First(entry).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}

	submission := new(EditForm)

	// Bind form data to struct
	if err := c.Bind(submission); err != nil {
		return c.String(http.StatusBadRequest, "Invalid submission")
	}

	// Validate submission
	if submission.UserText == "" {
		return c.String(http.StatusBadRequest, "Text cannot be empty")
	}

	if err = entries.Set(entry, submission.UserText, ollama.NewClientFromConfig()); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}

	referer := c.Request().Referer()
	if referer == "" {
		referer = "/history"
	}
	return c.Redirect(http.StatusSeeOther, referer)
}
