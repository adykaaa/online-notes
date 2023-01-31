package http

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog"
)

func PasetoAuth(c *PasetoCreator, symmetricKey string, l *zerolog.Logger) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tokenCookie, err := r.Cookie("authentication")
			if err != nil {
				http.Error(w, "auth cookie not set", http.StatusUnauthorized)
				l.Error().Err(err).Msgf("authentication cookie is not set!")
				return
			}

			payload, err := c.VerifyToken(tokenCookie.Value)
			if err != nil {
				if errors.Is(err, ErrTokenInvalid) {
					l.Error().Err(err).Msgf("PASETO is invalid!")
					http.Error(w, "invalid token", http.StatusUnauthorized)
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
