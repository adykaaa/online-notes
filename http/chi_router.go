package http

import (
	"github.com/adykaaa/online-notes/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewChiRouter(repo *db.Repository) *chi.Mux {
	router := chi.NewRouter()
	RegisterChiMiddlewares(router)
	RegisterChiHandlers(router, repo)

	return router
}

// TODO: set strict CORS when everything's gucci
func RegisterChiMiddlewares(r *chi.Mux) {
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Bearer"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

func RegisterChiHandlers(router *chi.Mux, repo *db.Repository) {
	router.Get("/", Home(repo))
	router.Post("/register", RegisterUser(repo))
}
