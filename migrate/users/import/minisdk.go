package importusers

import (
	util "github.com/mainflux/migrations/internal"
)

const (
	retrieveUserOps = "retrieving users"
)

// GetUserID returns the user ID associated with the email address provided.
func GetUserID(filePath, email string) string {
	if email == "" {
		return ""
	}
	records, err := util.ReadAllData(filePath, retrieveUserOps)
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
