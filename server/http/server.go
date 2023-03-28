package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

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
