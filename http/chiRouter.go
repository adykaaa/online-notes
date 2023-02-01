package http

import (
	"context"
	"net/http"
	"time"

	"github.com/adykaaa/httplog"
	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

func NewChiRouter(q sqlc.Querier, symmetricKey string, logger *zerolog.Logger) (*chi.Mux, error) {
	tokenCreator, err := NewPasetoCreator(symmetricKey)
	if err != nil {
		logger.Err(err).Msgf("could not create a new PasetoCreator. %v", err)
		return nil, err
	}

	router := chi.NewRouter()
	RegisterChiMiddlewares(router, logger)
	RegisterChiHandlers(router, q, tokenCreator, symmetricKey, logger)

	return router, nil
}

// TODO: set strict CORS when everything's gucci
func RegisterChiMiddlewares(r *chi.Mux, logger *zerolog.Logger) {
	// Request logger has middleware.Recoverer and RequestID baked into it.
	r.Use(render.SetContentType(render.ContentTypeJSON),
		httplog.RequestLogger(logger),
		middleware.Heartbeat("/ping"),
		middleware.RedirectSlashes,
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Bearer"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
}

func SetupHandler(w http.ResponseWriter, ctx context.Context) (*zerolog.Logger, context.Context, context.CancelFunc) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	l := zerolog.Ctx(ctx)
	return l, ctx, cancel
}

func RegisterChiHandlers(router *chi.Mux, q sqlc.Querier, c *PasetoCreator, symmetricKey string, logger *zerolog.Logger) {
	router.Post("/register", RegisterUser(q))
	router.Post("/login", LoginUser(q, c))
	router.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(c, symmetricKey, logger))
	})
}
