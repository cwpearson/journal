package handlers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"

	"github.com/cwpearson/journal/database"
	"github.com/cwpearson/journal/models"
)

type EntryGroup struct {
	Year    int
	Month   int
	Day     int
	Entries []models.Entry
}

func HistoryGet(c echo.Context) error {
	db := database.Get()

	var entries []models.Entry

	// Retrieve all entries ordered by date and N
	result := db.Preload("Tags").
		Order("year DESC, month DESC, day DESC, n ASC").
		Find(&entries)

	if result.Error != nil {
		return c.String(http.StatusInternalServerError, "Database error")
	}

	// Group entries by date
	groups := make(map[string]EntryGroup)
	for _, entry := range entries {
		key := fmt.Sprintf("%04d-%02d-%02d", entry.Year, entry.Month, entry.Day)
		group, exists := groups[key]
		if !exists {
			group = EntryGroup{
				Year:  entry.Year,
				Month: entry.Month,
				Day:   entry.Day,
			}
		}
		group.Entries = append(group.Entries, entry)
		groups[key] = group
	}

	// Convert map to sorted slice
	var sortedGroups []EntryGroup
	for _, group := range groups {
		sortedGroups = append(sortedGroups, group)
	}

	// Sort groups by date (newest first)
	sort.Slice(sortedGroups, func(i, j int) bool {
		a, b := sortedGroups[i], sortedGroups[j]
		if a.Year != b.Year {
			return a.Year > b.Year
		}
		if a.Month != b.Month {
			return a.Month > b.Month
		}
		return a.Day > b.Day
	})

	// Return HTML template
	return c.Render(http.StatusOK, "history.html", sortedGroups)
}
