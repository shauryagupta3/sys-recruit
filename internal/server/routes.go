package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type UserIdType float64

const UserID UserIdType = 0

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// default
	r.Get("/", s.HelloWorldHandler)
	r.Get("/health", s.healthHandler)

	//public
	r.Post("/signup", Make(s.handleSignup))
	r.Post("/login", Make(s.handleLogin))
	r.Get("/jobs", Make(s.GetAllJobs))
	r.Get("/jobs/{id}", Make(s.GetJobById))

	//protected applicant
	r.Group(func(r chi.Router) {
		r.Use(ApplicantOnly)
		r.Post("/uploadresume", Make(s.HandleUploadResume))
		r.Post("/jobs/{id}/apply", Make(s.ApplyToJobByID))
		r.Get("/myapplications",Make(s.GetJobsAppliedBy))
	})

	// protected routes admin
	r.Mount("/admin", s.adminRouter())

	return r
}

func (s *Server) adminRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(AdminOnly)
	r.Get("/", Make(s.handleAdmin))
	r.Post("/jobs", Make(s.handlePostJob))
	r.Get("/jobs", Make(s.handleGetJobsByAdmin)) // user not returned well not needed
	r.Get("/jobs/{id}",Make(s.handleGetJobsByIDAdmin)) 	// r.Get("/applicants", Make(s.handlePostJob)) all applicants

	return r
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id, err := AdminProtected(w, r); err == nil {
			ctx := context.WithValue(r.Context(), UserID, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
	})
}

func ApplicantOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id, err := ApplicantProtected(w, r); err == nil {
			ctx := context.WithValue(r.Context(), UserID, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
	})
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
