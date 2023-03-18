package httplib

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Msg map[string]string

func JSON(w http.ResponseWriter, payload interface{}, code int) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error while marshaling for the response"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func SetupHandler(w http.ResponseWriter, ctx context.Context) (*zerolog.Logger, context.Context, context.CancelFunc) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	l := zerolog.Ctx(ctx)
	return l, ctx, cancel
}

func SetCookie(w http.ResponseWriter, name string, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    token,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   true,
	})
}
