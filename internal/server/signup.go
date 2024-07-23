package server

import (
	"encoding/json"
	"net/http"
	"recruit-sys/internal/models"
)

func (s *Server)handleSignup(w http.ResponseWriter, r *http.Request) error {

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return InvalidJson()
	}

	if err := user.Validate(); err != nil {
		return InvalidReqJsonData(err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
	return nil
}
