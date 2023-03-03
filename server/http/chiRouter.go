package server

import (
	"time"

	"github.com/adykaaa/httplog"
	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	auth "github.com/adykaaa/online-notes/server/http/auth"
	"github.com/adykaaa/online-notes/server/http/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

func RegisterChiMiddlewares(r *chi.Mux, l *zerolog.Logger) {
	// Request logger has middleware.Recoverer and RequestID baked into it.
	r.Use(render.SetContentType(render.ContentTypeJSON),
		httplog.RequestLogger(l),
		middleware.Heartbeat("/ping"),
		middleware.RedirectSlashes,
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost:3000"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Bearer", "Set-Cookie", "X-Powered-By", "X-Content-Type-Options"},
			ExposedHeaders:   []string{"Link", "Access-Control-Expose-Headers"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
}

func RegisterChiHandlers(router *chi.Mux, q sqlc.Querier, t auth.TokenManager, tokenDuration time.Duration, l *zerolog.Logger) {
	router.Post("/register", handlers.RegisterUser(q))
	router.Post("/login", handlers.LoginUser(q, t, tokenDuration))
	router.Post("/logout", handlers.LogoutUser())
	router.Route("/notes", func(router chi.Router) {
		router.Use(auth.AuthMiddleware(t, l))
		router.Post("/create", handlers.CreateNote(q))
		router.Get("/", handlers.GetAllNotesFromUser(q))
		router.Put("/{id}", handlers.UpdateNote(q))
		router.Delete("/{id}", handlers.DeleteNote(q))
	})
}

func NewChiRouter(q sqlc.Querier, symmetricKey string, tokenDuration time.Duration, l *zerolog.Logger) (*chi.Mux, error) {
	tokenCreator, err := auth.NewPasetoManager(symmetricKey)
	if err != nil {
		l.Err(err).Msgf("could not create a new PasetoCreator. %v", err)
		return nil, err
	}
	router := chi.NewRouter()
	RegisterChiMiddlewares(router, l)
	RegisterChiHandlers(router, q, tokenCreator, tokenDuration, l)

	return router, nil
}