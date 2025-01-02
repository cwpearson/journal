package entries

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"

	"github.com/cwpearson/journal/config"
	"github.com/cwpearson/journal/database"
	"github.com/cwpearson/journal/models"
	"github.com/cwpearson/journal/ollama"
	"github.com/cwpearson/journal/tags"
)

func GetForDate(year, month, day int) ([]models.Entry, error) {
	db := database.Get()
	res := []models.Entry{}
	err := db.Where("year = ?", year).Where("month = ?", month).Where("day = ?", day).Find(&res).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return res, nil
}

func Create(client *ollama.Client, text string) (*models.Entry, error) {
	db := database.Get()

	now := time.Now()

	year := now.Year()        // Returns year as int
	month := int(now.Month()) // Month() returns time.Month, so we cast to int
	day := now.Day()          // Returns day as int

	entryDir := filepath.Join(config.DataDir(),
		fmt.Sprintf("%04d", year),
		fmt.Sprintf("%02d", month),
		fmt.Sprintf("%02d", day),
	)

	log.Println("create", entryDir)
	err := os.MkdirAll(entryDir, 0755)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	// find the largest N
	others, err := GetForDate(year, month, day)
	if err != nil {
		return nil, err
	}
	nextN := 0
	if len(others) > 0 {
		for _, e := range others {
			if e.N > nextN {
				nextN = e.N
			}
		}
		nextN++
	}

	textPath := filepath.Join(entryDir, fmt.Sprintf("%03d.txt", nextN))
	log.Println("write", textPath)
	err = os.WriteFile(textPath, []byte(text), 0644)
	if err != nil {
		return nil, err
	}

	database.Lock()
	defer database.Unlock()
	entry := &models.Entry{
		Year:  year,
		Month: month,
		Day:   day,
		N:     nextN,
	}
	err = db.Create(&entry).Error
	if err != nil {
		os.Remove(textPath) // ignore errors
		return nil, err
	}

	// generate and add tags
	go func() {
		kwds, err := client.Keywords(text)
		if err != nil {
			log.Println("keywords error:", err)
		}

		SetKeywords(entry, kwds)
	}()

	// generate and add summary
	go func() {
		summary, err := client.Summary(text)
		if err != nil {
			log.Println("summary error:", err)
		}
		if summary == "" {
			summary = text
		}
		if err := SetSummary(entry, summary); err != nil {
			log.Println("SetSummary error:", err)
		}
	}()

	return entry, nil
}

func AddKeyword(entry *models.Entry, keyword string) error {
	keyword = tags.Clean(keyword)

	db := database.Get()
	database.Lock()
	defer database.Unlock()
	log.Println("Add keyword", keyword, "to entry", entry.ID)

	// First try to find existing tag
	var tag models.Tag
	result := db.Where("s = ?", keyword).First(&tag)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Tag doesn't exist, create new one
			tag = models.Tag{
				S: keyword,
			}
			if err := db.Create(&tag).Error; err != nil {
				return fmt.Errorf("failed to create tag: %w", err)
			}
		} else {
			return fmt.Errorf("error finding tag: %w", result.Error)
		}
	}

	// Associate tag with entry
	if err := db.Model(&entry).Association("Tags").Append(&tag); err != nil {
		return fmt.Errorf("failed to associate tag with entry: %w", err)
	}

	return nil
}

func SetSummary(entry *models.Entry, summary string) error {
	entry.Summary = summary
	database.Lock()
	defer database.Unlock()
	log.Println("Set summary for entry", entry.ID)
	return database.Get().Save(entry).Error
}
