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

func RegisterChiHandlers(r *chi.Mux, s NoteService, t auth.TokenManager, td time.Duration, l *zerolog.Logger) {
	r.Post("/register", RegisterUser(s))
	r.Post("/login", LoginUser(s, t, td))
	r.Post("/logout", LogoutUser())
	r.Route("/notes", func(r chi.Router) {
		r.Use(auth.AuthMiddleware(t, l))
		r.Post("/create", CreateNote(s))
		r.Get("/", GetAllNotesFromUser(s))
		r.Put("/{id}", UpdateNote(s))
		r.Delete("/{id}", DeleteNote(s))
	})
}

func NewChiRouter(s NoteService, symmetricKey string, tokenDuration time.Duration, l *zerolog.Logger) (*chi.Mux, error) {
	pm, err := auth.NewPasetoManager(symmetricKey)
	if err != nil {
		l.Err(err).Msgf("could not create a new PasetoCreator. %v", err)
		return nil, err
	}

	r := chi.NewRouter()
	RegisterChiMiddlewares(r, l)
	RegisterChiHandlers(r, s, pm, tokenDuration, l)

	return r, nil
}
