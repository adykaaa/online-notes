package server

import (
	"context"
	"time"

	"github.com/adykaaa/httplog"
	db "github.com/adykaaa/online-notes/db/sqlc"
	auth "github.com/adykaaa/online-notes/server/http/auth"
	handlers "github.com/adykaaa/online-notes/server/http/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type NoteService interface {
	CreateNote(ctx context.Context, title string, username string, text string) (uuid.UUID, error)
	GetAllNotesFromUser(ctx context.Context, username string) ([]db.Note, error)
	DeleteNote(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	UpdateNote(ctx context.Context, reqID uuid.UUID, title string, text string, isTextEmpty bool) (uuid.UUID, error)
	RegisterUser(ctx context.Context, args *db.RegisterUserParams) (string, error)
	GetUser(ctx context.Context, username string) (db.User, error)
}

func registerChiMiddlewares(r *chi.Mux, l *zerolog.Logger) {
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

func registerChiHandlers(r *chi.Mux, s NoteService, t auth.TokenManager, tokenDuration time.Duration, l *zerolog.Logger) {
	r.Post("/register", handlers.RegisterUser(s))
	r.Post("/login", handlers.LoginUser(s, t, tokenDuration))
	r.Post("/logout", handlers.LogoutUser())
	r.Route("/notes", func(r chi.Router) {
		r.Use(auth.AuthMiddleware(t, l))
		r.Post("/create", handlers.CreateNote(s))
		r.Get("/", handlers.GetAllNotesFromUser(s))
		r.Put("/{id}", handlers.UpdateNote(s))
		r.Delete("/{id}", handlers.DeleteNote(s))
	})
}

func NewChiRouter(s NoteService, symmetricKey string, tokenDuration time.Duration, l *zerolog.Logger) (*chi.Mux, error) {
	pm, err := auth.NewPasetoManager(symmetricKey)
	if err != nil {
		l.Err(err).Msgf("could not create a new PasetoCreator. %v", err)
		return nil, err
	}

	r := chi.NewRouter()
	registerChiMiddlewares(r, l)
	registerChiHandlers(r, s, pm, tokenDuration, l)

	return r, nil
}
