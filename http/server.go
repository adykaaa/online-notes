package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

type Server struct {
	logger          *zerolog.Logger
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func NewServer(handler http.Handler, addr string, logger *zerolog.Logger) (*Server, error) {

	s := &Server{
		server: &http.Server{
			Handler: handler,
			Addr:    addr,
		},
		notify:          make(chan error, 1),
		shutdownTimeout: 5 * time.Second,
		logger:          logger,
	}

	if addr == "" {
		s.logger.Error().Msg("Server address cannot be empty!")
	}

	s.Start()
	return s, nil
}

func (s *Server) Start() {
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
		s.logger.Info().Msgf("Server connection error %v", err)
	}
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
