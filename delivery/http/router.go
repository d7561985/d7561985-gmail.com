package http

import (
	"github.com/go-chi/chi"
	c "github.com/go-chi/cors"
)

const (
	h         = "/health"
	questions = "/questions"
)

func (s *service) Router() chi.Router {
	r := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := c.New(c.Options{
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r.Use(cors.Handler)
	r.Use(s.Logger)
	r.Use(s.Jaeger)
	r.Use(s.Recovery)

	r.Get(h, s.HealthHandler)
	r.Mount(questions, s.route(questions))

	return r
}
