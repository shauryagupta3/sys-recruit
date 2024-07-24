package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	userID, err := GetUserIDFromJWT(tokenString)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userID)
	return nil
}
