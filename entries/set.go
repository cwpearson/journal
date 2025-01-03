package entries

import (
	"fmt"
	"log"
	"os"

	"github.com/cwpearson/journal/database"
	"github.com/cwpearson/journal/models"
	"github.com/cwpearson/journal/ollama"
)

func Set(entry *models.Entry, text string, client *ollama.Client) error {
	textPath := TextPath(entry)
	log.Println("write", textPath)
	err := os.WriteFile(textPath, []byte(text), 0644)
	if err != nil {
		return err
	}

	database.Lock()
	defer database.Unlock()

	return Analyze(entry, client)
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
