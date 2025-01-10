package ollama

import (
	"time"

	"gorm.io/gorm"

	"github.com/cwpearson/journal/database"
)

// a record of ollama API calls and their results
type Record struct {
	gorm.Model
	When   time.Time
	Kind   string
	Result string
}

func createRecord(kind, result string) {
	database.Get().Create(&Record{
		When:   time.Now(),
		Kind:   kind,
		Result: result,
	})
}
