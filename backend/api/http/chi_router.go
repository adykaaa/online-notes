package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// TODO: move this interface away from here where it's used
type Router interface {
	Get(path string, handlerFunc func(w http.ResponseWriter, r *http.Request))
	Post(path string, handlerFunc func(w http.ResponseWriter, r *http.Request))
	Delete(path string, handlerFunc func(w http.ResponseWriter, r *http.Request))
}

func RegisterChiMiddlewares(r *chi.Mux) {
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Bearer"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
}

func RegisterChiRoutes(r *chi.Mux) {
}
