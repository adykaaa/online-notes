package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type TokenManager interface {
	CreateToken(username string, duration time.Duration) (string, *PasetoPayload, error)
	VerifyToken(token string) (*PasetoPayload, error)
}

func AuthMiddleware(t TokenManager, l *zerolog.Logger) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tokenCookie, err := r.Cookie("paseto")
			if err != nil {
				http.Error(w, "paseto auth cookie not set", http.StatusUnauthorized)
				l.Error().Err(err).Msgf("authentication cookie is not set!")
				return
			}

			payload, err := t.VerifyToken(tokenCookie.Value)
			if err != nil {
				if errors.Is(err, ErrTokenInvalid) {
					l.Error().Err(err).Msgf("PASETO is invalid!")
					http.Error(w, "invalid token", http.StatusForbidden)
					return
				}
				if errors.Is(err, ErrTokenExpired) {
					l.Error().Err(err).Msgf("PASETO has expired")
					http.Error(w, "token has expired", http.StatusUnauthorized)
					return
				}
			}

			l.Info().Msgf("User %s is authorized! tokenID: %v, issuedAt: %v, expiresAt: %v", payload.Username, payload.ID, payload.IssuedAt, payload.ExpiresAt)
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return f
}
