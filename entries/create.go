package entries

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/gorm"

	"github.com/cwpearson/journal/database"
	"github.com/cwpearson/journal/models"
	"github.com/cwpearson/journal/ollama"
	"github.com/cwpearson/journal/tags"
)

func Create(year, month, day int, text string) (*models.Entry, error) {
	db := database.Get()
	database.Lock()
	defer database.Unlock()

	entry := &models.Entry{
		Year:  year,
		Month: month,
		Day:   day,
	}
	err := db.Create(&entry).Error
	if err != nil {
		return nil, err
	}

	textPath := TextPath(entry)
	entryDir := filepath.Dir(textPath)
	log.Println("create", entryDir)
	err = os.MkdirAll(entryDir, 0755)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	log.Println("write", textPath)
	err = os.WriteFile(textPath, []byte(text), 0644)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func Analyze(entry *models.Entry, client *ollama.Client) error {

	// load text
	bytes, err := os.ReadFile(TextPath(entry))
	if err != nil {
		return err
	}
	text := string(bytes)

	// generate and add tags
	go func() {

		// clear keywords
		err = SetKeywords(entry, []string{})
		if err != nil {
			log.Println("SetKeywords error:", err)
		}

		kwds, err := client.Keywords(text)
		if err != nil {
			log.Println("keywords error:", err)
		}

		SetKeywords(entry, kwds)
	}()

	// generate and add summary
	go func() {
		// clear summary
		if err := SetSummary(entry, ""); err != nil {
			log.Println("SetSummary error:", err)
		}

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

	return nil
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
