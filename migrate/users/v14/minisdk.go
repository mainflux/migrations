package users14

import (
	"github.com/mainflux/migrations/internal/util"
)

// GetUserID returns the user ID associated with the email address provided
func GetUserID(filePath, email string) string {
	if email == "" {
		return ""
	}
	records, err := util.ReadData(filePath)
	if err != nil {
		return ""
	}
	for _, record := range records {
		if record[1] == email {
			return record[0]
		}
	}
	return email
}
