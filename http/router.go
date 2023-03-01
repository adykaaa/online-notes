package http

import (
	"net/http"
	"time"
)

type Router interface {
	Get(pattern string, handlerFn http.HandlerFunc)
	Delete(pattern string, handlerFn http.HandlerFunc)
	Post(pattern string, handlerFn http.HandlerFunc)
	Put(pattern string, handlerFn http.HandlerFunc)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type TokenManager interface {
	CreateToken(username string, duration time.Duration) (string, *PasetoPayload, error)
	VerifyToken(token string) (*PasetoPayload, error)
}
