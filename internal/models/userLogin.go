package models

import (
	"regexp"
)

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserLogin) Validate() map[string]string {
	errors := make(map[string]string)

	if !isValidEmail(u.Email) {
		errors["email"] = ("invalid email address")
	}

	if len(u.Password) < 3 {
		errors["password"] = ("password must be at least 3 characters long")
	}

	return nil
}

// Helper function to validate email format
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
