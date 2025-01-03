package entries

import (
	"fmt"
	"path/filepath"

	"github.com/cwpearson/journal/config"
	"github.com/cwpearson/journal/models"
)

func TextPath(e *models.Entry) string {
	return filepath.Join(
		config.DataDir(),
		fmt.Sprintf("%04d", e.Year),
		fmt.Sprintf("%02d", e.Month),
		fmt.Sprintf("%02d.txt", e.Day),
	)
}
