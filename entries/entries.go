package entries

import (
	"fmt"

	"github.com/cwpearson/journal/models"
)

func PathComponents(e *models.Entry) (string, string, string, string) {
	return fmt.Sprintf("%04d", e.Year), fmt.Sprintf("%02d", e.Month), fmt.Sprintf("%02d", e.Day), fmt.Sprintf("%03d.txt", e.N)
}
