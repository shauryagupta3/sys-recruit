package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type APIError struct {
	StatusCode int `json:"status_code"`
	Message    any `json:"message"`
}

func InvalidReqJsonData(errors map[string]string) APIError {
	return APIError{
		StatusCode: http.StatusBadRequest,
		Message:    errors,
	}
}

func InvalidJson() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid json data"))
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %d", e.StatusCode)
}

func NewAPIError(code int, msg error) APIError {
	return APIError{
		StatusCode: code,
		Message:    msg.Error(),
	}
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func Make(h APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if apiError, ok := err.(APIError); ok {
				WriteJSON(w, apiError.StatusCode, apiError)
			} else {
				WriteJSON(w, http.StatusInternalServerError, map[string]any{"statusCode": http.StatusInternalServerError, "message": err})
			}
			slog.Error("API err", "err", err.Error(), "path", r.URL.Path)
		}
	}
}

func WriteJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}
