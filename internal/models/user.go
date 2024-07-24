package models

import (
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

