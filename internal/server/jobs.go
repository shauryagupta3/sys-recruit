package server

import (
	"encoding/json"
	"net/http"
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
