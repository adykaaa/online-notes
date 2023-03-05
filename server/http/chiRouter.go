package server

import (
	"context"
	"time"

	"github.com/adykaaa/httplog"
	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	auth "github.com/adykaaa/online-notes/server/http/auth"
	handler "github.com/adykaaa/online-notes/server/http/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type NoteService interface {
	CreateNote(ctx context.Context, title string, username string, text string) (uuid.UUID, error)
	GetAllNotesFromUser(ctx context.Context, username string) ([]sqlc.Note, error)
	DeleteNote(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	UpdateNote(ctx context.Context, reqID uuid.UUID, title string, text string, isTextEmpty bool) (uuid.UUID, error)
	RegisterUser(ctx context.Context, username string, password string, email string) (string, error)
	GetUser(ctx context.Context, username string) (sqlc.User, error)
}

func RegisterChiMiddlewares(r *chi.Mux, l *zerolog.Logger) {
	// Request logger has middleware.Recoverer and RequestID baked into it.
	r.Use(httplog.RequestLogger(l),
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

func RegisterChiHandlers(router *chi.Mux, s NoteService, t auth.TokenManager, tokenDuration time.Duration, l *zerolog.Logger) {
	router.Post("/register", handler.RegisterUser(s))
	router.Post("/login", handler.LoginUser(s, t, tokenDuration))
	router.Post("/logout", handler.LogoutUser())
	router.Route("/notes", func(router chi.Router) {
		router.Use(auth.AuthMiddleware(t, l))
		router.Post("/create", handler.CreateNote(s))
		router.Get("/", handler.GetAllNotesFromUser(s))
		router.Put("/{id}", handler.UpdateNote(s))
		router.Delete("/{id}", handler.DeleteNote(s))
	})
}

func NewChiRouter(s NoteService, symmetricKey string, tokenDuration time.Duration, l *zerolog.Logger) (*chi.Mux, error) {
	tokenManager, err := auth.NewPasetoManager(symmetricKey)
	if err != nil {
		l.Err(err).Msgf("could not create a new PasetoCreator. %v", err)
		return nil, err
	}
	router := chi.NewRouter()
	RegisterChiMiddlewares(router, l)
	RegisterChiHandlers(router, s, tokenManager, tokenDuration, l)

	return router, nil
}
