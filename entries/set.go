package entries

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cwpearson/journal/config"
	"github.com/cwpearson/journal/database"
	"github.com/cwpearson/journal/models"
	"github.com/cwpearson/journal/ollama"
)

func Set(entryId uint, text string, client *ollama.Client) error {
	db := database.Get()
	entry := new(models.Entry)
	err := db.Where("id = ?", entryId).Find(entry).Error
	if err != nil {
		return err
	}

	year, month, day, name := PathComponents(entry)

	textPath := filepath.Join(config.DataDir(), year, month, day, name)
	log.Println("write", textPath)
	err = os.WriteFile(textPath, []byte(text), 0644)
	if err != nil {
		return err
	}

	database.Lock()
	defer database.Unlock()

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

		err = SetKeywords(entry, kwds)
		if err != nil {
			log.Println("SetKeywords error:", err)
		}
	}()

	// generate and add summary
	go func() {
		// clear summary
		if err := SetSummary(entry, ""); err != nil {
			log.Println("SetSummary error:", err)
		}

		// generate summary
		summary, err := client.Summary(text)
		if err != nil {
			log.Println("summary error:", err)
		}

		// set generated summary
		if summary == "" {
			summary = text
		}
		if err := SetSummary(entry, summary); err != nil {
			log.Println("SetSummary error:", err)
		}
	}()

	return nil
}

func SetKeywords(entry *models.Entry, keywords []string) error {

	db := database.Get()

	// clear tags
	err := db.Model(&entry).Association("Tags").Clear()
	if err != nil {
		return err
	}

	for _, kwd := range keywords {
		err := AddKeyword(entry, kwd)
		if err != nil {
			fmt.Println("AddKeyword error:", err)
		}
	}

	return nil
}
