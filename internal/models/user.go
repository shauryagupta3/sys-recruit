package models

import (
	"strings"
	"time"
)

type User struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Address         string `json:"address"`
	UserType        string `json:"user_type"`
	PasswordHash    string `json:"password_hash"`
	ProfileHeadline string `json:"profile_headline"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) Validate() map[string]string {
	errors := make(map[string]string)

	if len(u.Name) < 3 {
		errors["name"] = ("name at least 3 chars")
	}

	if u.UserType != "admin" && u.UserType != "applicant" {
		errors["user_type"] = "invalid user type"
	}

	if !isValidEmail(u.Email) {
		errors["email"] = ("invalid email address")
	}

	if len(u.PasswordHash) < 3 {
		errors["password"] = ("password must be at least 3 characters long")
	}

	return errors
}

func ProcessUserInput(user *User) {
	user.Name = strings.ToLower(strings.TrimSpace(user.Name))
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Address = strings.ToLower(strings.TrimSpace(user.Address))
	user.UserType = strings.ToLower(strings.TrimSpace(user.UserType))
	user.ProfileHeadline = strings.ToLower(strings.TrimSpace(user.ProfileHeadline))
}
