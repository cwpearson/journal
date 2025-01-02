package database

import (
	"path/filepath"
	"sync"

	"github.com/cwpearson/journal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var mu sync.Mutex
var db *gorm.DB

func Init() error {
	dbPath := filepath.Join(config.ConfigDir(), "journal.sqlite")
	var err error
	db, err = gorm.Open(sqlite.Open(dbPath))
	if err != nil {
		return err
	}
	return nil
}

func Get() *gorm.DB {
	return db
}

func Lock() {
	mu.Lock()
}

func Unlock() {
	mu.Unlock()
}
