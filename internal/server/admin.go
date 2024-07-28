package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recruit-sys/internal/models"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) error {
	userID, ok := r.Context().Value(UserID).(float64)
	if !ok {
		return NewAPIError(http.StatusBadGateway, fmt.Errorf("unable to proceed"))
	}

	user, err := s.db.SelectUserWhereID(userID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Hello admin id : " + fmt.Sprintf("%d", user.ID))
	return nil
}

func (s *Server) handlePostJob(w http.ResponseWriter, r *http.Request) error {
	var job models.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		return InvalidJson()
	}
	userID, ok := r.Context().Value(UserID).(float64)
	if !ok {
		return NewAPIError(http.StatusBadGateway, fmt.Errorf("unable to proceed"))
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

func (s *Server) handleGetJobsByAdmin(w http.ResponseWriter, r *http.Request) error {
	userID, ok := r.Context().Value(UserID).(float64)
	if !ok {
		return NewAPIError(http.StatusBadGateway, fmt.Errorf("unable to proceed"))
	}

	jobs, err := s.db.SelectJobsPostedBy(userID)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jobs)
	return nil
}

func (s *Server) handleGetJobsByIDAdmin(w http.ResponseWriter, r *http.Request) error {
	userID, ok := r.Context().Value(UserID).(float64)
	if !ok {
		return NewAPIError(http.StatusBadGateway, fmt.Errorf("unable to proceed"))
	}

	jobId := chi.URLParam(r, "id")
	i, err := strconv.Atoi(jobId)
	if err != nil {
		return NewAPIError(http.StatusBadRequest, fmt.Errorf("bad id given"))
	}

	job, err := s.db.SelectJobByIdAdmin(userID, i)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(job)
	return nil
}
