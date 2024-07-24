package server

import (
	"encoding/json"
	"net/http"
	"recruit-sys/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) error {

	var userLogin models.UserLogin

	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		return InvalidJson()
	}

	if err := userLogin.Validate(); err != nil {
		return InvalidReqJsonData(err)
	}

	user, err := s.db.SelectUserWhereMail(userLogin.Email)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
	return nil
}
