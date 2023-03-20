package users

// GetUserID returns the user ID associated with the email address provided
func GetUserID(email string) string {
	// TODO Return UseID from UserEmail
	if email == "" {
		return ""
	}
	return email
}
