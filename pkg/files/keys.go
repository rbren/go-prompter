package files

import (
	"time"
)

func GetHeadlineKey() string {
	// get current date in YYYY-MM-DD format
	date := time.Now().Format("2006-01-02")
	return "headlines/" + date + ".json"
}
