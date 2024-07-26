package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (s *Server) GetAllJobs(w http.ResponseWriter, r *http.Request) error {
	jobs, err := s.db.SelectAllJobs()
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jobs)
	return nil
}

func (s *Server) GetJobById(w http.ResponseWriter, r *http.Request) error {
	jobId := chi.URLParam(r, "id")
	i, err := strconv.Atoi(jobId)
	if err != nil {
		return NewAPIError(http.StatusBadRequest, fmt.Errorf("bad id given"))
	}
	jobs, err := s.db.SelectJobsByID(i)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jobs)
	return nil
}
