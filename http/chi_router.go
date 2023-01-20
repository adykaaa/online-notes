package http

import (
	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
)

func NewChiRouter(q sqlc.Querier, logger zerolog.Logger) *chi.Mux {
	router := chi.NewRouter()
	RegisterChiMiddlewares(router, logger)
	RegisterChiHandlers(router, q)

	return router
}

// TODO: set strict CORS when everything's gucci
func RegisterChiMiddlewares(r *chi.Mux, logger zerolog.Logger) {

	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Bearer"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

func RegisterChiHandlers(router *chi.Mux, q sqlc.Querier) {
	router.Get("/", Home(q))
	router.Post("/register", RegisterUser(q))
}
