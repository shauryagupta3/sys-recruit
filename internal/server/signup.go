package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recruit-sys/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) error {

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return InvalidJson()
	}
	models.ProcessUserInput(&user)

	if err := user.Validate(); len(err)>0 {
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
