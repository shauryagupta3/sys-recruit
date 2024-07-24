package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recruit-sys/internal/models"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) error {

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return InvalidJson()
	}
	ProcessUserInput(&user)

	if err := signupValidate(&user); len(err) > 0 {
		fmt.Println(err)
		return InvalidReqJsonData(err)
	}
	PasswordHash, err := hashPassword(user.PasswordHash)
	if err != nil {
		return err
	}
	user.PasswordHash = PasswordHash

	err = s.db.CreateUser(&user)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
	return nil
}

func hashPassword(password string) (string, error) {
	pass := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func signupValidate(u *models.User) map[string]string {
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

func ProcessUserInput(user *models.User) {
	user.Name = strings.ToLower(strings.TrimSpace(user.Name))
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Address = strings.ToLower(strings.TrimSpace(user.Address))
	user.UserType = strings.ToLower(strings.TrimSpace(user.UserType))
	user.ProfileHeadline = strings.ToLower(strings.TrimSpace(user.ProfileHeadline))
}

// Helper function to validate email format
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
