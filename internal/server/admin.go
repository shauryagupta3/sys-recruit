package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recruit-sys/internal/models"
)

func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	userID, err := GetUserIDFromJWT(tokenString)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Hello admin id : " + fmt.Sprintf("%f", userID))
	return nil
}

func (s *Server) handlePostJob(w http.ResponseWriter, r *http.Request) error {
	var job models.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		return InvalidJson()
	}
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	userID, err := GetUserIDFromJWT(tokenString)
	if err != nil {
		return err
	}
	user, err := s.db.SelectUserWhereID(userID)
	if err != nil {
		return err
	}

	job.PostedBy = user
	job.PostedByID = user.ID

	err = s.db.CreateJob(&job)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(job)
	return nil
}
