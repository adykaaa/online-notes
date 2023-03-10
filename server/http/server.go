package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type NoteService interface {
	CreateNote(ctx context.Context, title string, username string, text string) (uuid.UUID, error)
	GetAllNotesFromUser(ctx context.Context, username string) ([]sqlc.Note, error)
	DeleteNote(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	UpdateNote(ctx context.Context, reqID uuid.UUID, title string, text string, isTextEmpty bool) (uuid.UUID, error)
	RegisterUser(ctx context.Context, username string, password string, email string) (string, error)
	GetUser(ctx context.Context, username string) (sqlc.User, error)
}

type Router interface {
	Get(pattern string, handlerFn http.HandlerFunc)
	Delete(pattern string, handlerFn http.HandlerFunc)
	Post(pattern string, handlerFn http.HandlerFunc)
	Put(pattern string, handlerFn http.HandlerFunc)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type HTTPServer struct {
	logger          *zerolog.Logger
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func NewHTTP(r Router, addr string, l *zerolog.Logger) (*HTTPServer, error) {
	s := &HTTPServer{
		server: &http.Server{
			Handler: r,
			Addr:    addr,
		},
		notify:          make(chan error, 1),
		shutdownTimeout: 5 * time.Second,
		logger:          l,
	}

	if addr == "" {
		s.logger.Error().Msg("server address is empty")
		return nil, errors.New("server address cannot be empty")
	}

	s.Start()
	return s, nil
}

func (s *HTTPServer) Start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		s.logger.Info().Msgf("Server run interrupted by OS signal %s", sig.String())
	case err := <-s.Notify():
		s.logger.Error().Msgf("error during server connection %v", err)
	}
}

func (s *HTTPServer) Notify() <-chan error {
	return s.notify
}

func (s *HTTPServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
