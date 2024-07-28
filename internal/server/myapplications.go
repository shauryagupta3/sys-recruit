package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) GetJobsAppliedBy(w http.ResponseWriter, r *http.Request) error {

	userID, ok := r.Context().Value(UserID).(float64)
	if !ok {
		return NewAPIError(http.StatusBadGateway, fmt.Errorf("unable to proceed"))
	}
	jobs, err := s.db.SelectJobsAppliedBy(userID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jobs)
	return nil
}
