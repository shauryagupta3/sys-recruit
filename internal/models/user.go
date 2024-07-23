package models

type User struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Address         string `json:"address"`
	UserType        string `json:"user_type"`
	PasswordHash    string `json:"password_hash"`
	ProfileHeadline string `json:"profile_headline"`
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
