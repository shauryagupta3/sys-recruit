package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// default
	r.Get("/", s.HelloWorldHandler)
	r.Get("/health", s.healthHandler)

	//public
	r.Post("/signup", Make(s.handleSignup))
	r.Post("/login", Make(s.handleLogin))
	r.Get("/jobs",Make(s.GetAllJobs))

	//protected applicant
	r.Group(func(r chi.Router) {
		r.Use(ApplicantOnly)
		r.Post("/uploadresume", Make(s.HandleUploadResume))
	})

	// protected routes admin
	r.Mount("/admin", s.adminRouter())

	return r
}

func (s *Server) adminRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(AdminOnly)
	r.Get("/", Make(s.handleAdmin))
	r.Post("/job", Make(s.handlePostJob))

	return r
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := AdminProtected(w, r); err == nil {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
	})
}

func ApplicantOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := ApplicantProtected(w, r); err == nil {
			next.ServeHTTP(w, r)
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
