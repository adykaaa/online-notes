package http

import (
	"time"

	"github.com/adykaaa/httplog"
	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

type TokenManager interface {
	CreateToken(username string, duration time.Duration) (string, *PasetoPayload, error)
	VerifyToken(token string) (*PasetoPayload, error)
}

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

type msg map[string]string

func RegisterChiHandlers(router *chi.Mux, q sqlc.Querier, t TokenManager, tokenDuration time.Duration, l *zerolog.Logger) {
	router.Post("/register", RegisterUser(q))
	router.Post("/login", LoginUser(q, t, tokenDuration))
	router.Post("/logout", LogoutUser())
	router.Route("/notes", func(router chi.Router) {
		router.Use(AuthMiddleware(t, l))
		router.Post("/create", CreateNote(q))
		router.Get("/", GetAllNotesFromUser(q))
		router.Put("/{id}", UpdateNote(q))
		router.Delete("/{id}", DeleteNote(q))
	})
}

func NewChiRouter(q sqlc.Querier, symmetricKey string, tokenDuration time.Duration, l *zerolog.Logger) (*chi.Mux, error) {
	tokenCreator, err := NewPasetoCreator(symmetricKey)
	if err != nil {
		l.Err(err).Msgf("could not create a new PasetoCreator. %v", err)
		return nil, err
	}
	router := chi.NewRouter()
	RegisterChiMiddlewares(router, l)
	RegisterChiHandlers(router, q, tokenCreator, tokenDuration, l)

	return router, nil
}
