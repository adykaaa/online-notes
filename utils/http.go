package utils

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

func JSON(w http.ResponseWriter, payload interface{}, code int) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error while waiting for the response"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func SetupHandler(w http.ResponseWriter, ctx context.Context) (*zerolog.Logger, context.Context, context.CancelFunc) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	l := zerolog.Ctx(ctx)
	return l, ctx, cancel
}
