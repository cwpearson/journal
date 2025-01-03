package entries

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/gorm"

	"github.com/cwpearson/journal/database"
	"github.com/cwpearson/journal/models"
)

func removeIfEmpty(dirPath string) error {
	// Read directory contents
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// Remove if empty
	if len(entries) == 0 {
		log.Println("remove empty directory", dirPath)
		if err := os.Remove(dirPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove directory: %w", err)
		}
	}

	return nil
}
func Delete(entryId uint) error {

	db := database.Get()
	entry := new(models.Entry)

	database.Lock()
	defer database.Unlock()

	err := db.Where("id = ?", entryId).Find(entry).Error
	if err != nil {
		return err
	}

	// erase file
	textPath := TextPath(entry)
	log.Println("remove", textPath)
	err = os.Remove(textPath)
	if err != nil {
		return err
	}

	monthDir := filepath.Dir(textPath)
	err = removeIfEmpty(monthDir)
	if err != nil {
		return err
	}
	yearDir := filepath.Dir(monthDir)
	err = removeIfEmpty(yearDir)
	if err != nil {
		return err
	}

	// clear tags
	err = db.Model(&entry).Association("Tags").Clear()
	if err != nil {
		return err
	}

	// delete all tags with no entries
	err = GarbageCollectTags(db)
	if err != nil {
		return err
	}

	// delete entry record
	return db.Model(&models.Entry{}).Delete("id = ?", entryId).Error
}

func GarbageCollectTags(db *gorm.DB) error {
	return db.Where("id NOT IN (SELECT DISTINCT tag_id FROM entry_tags)").Delete(&models.Tag{}).Error
}
