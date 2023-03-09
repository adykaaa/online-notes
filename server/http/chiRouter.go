package server

import (
	"time"

	"github.com/adykaaa/httplog"
	auth "github.com/adykaaa/online-notes/server/http/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
)

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
	router.Post("/register", RegisterUser(s))
	router.Post("/login", LoginUser(s, t, tokenDuration))
	router.Post("/logout", LogoutUser())
	router.Route("/notes", func(router chi.Router) {
		router.Use(auth.AuthMiddleware(t, l))
		router.Post("/create", CreateNote(s))
		router.Get("/", GetAllNotesFromUser(s))
		router.Put("/{id}", UpdateNote(s))
		router.Delete("/{id}", DeleteNote(s))
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
