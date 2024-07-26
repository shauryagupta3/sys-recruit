package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (s *Server) ApplyToJobByID(w http.ResponseWriter, r *http.Request) error {
	userID, ok := r.Context().Value(UserID).(float64)
	if !ok {
		return NewAPIError(http.StatusBadGateway, fmt.Errorf("unable to proceed"))
	}
	JobID := chi.URLParam(r, "id")
	i, err := strconv.Atoi(JobID)
	if err != nil {
		return err
	}
	err = s.db.ApplyToJob(i, int(userID))
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fmt.Sprintf("success applied to job : %d", i))
	return nil
}
