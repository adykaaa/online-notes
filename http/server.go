package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func NewServer(handler http.Handler, addr string) (*Server, error) {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	if addr == "" {
		fmt.Errorf("Address of the server must be specified!")
	}

	s := &Server{
		server: &http.Server{
			Handler: handler,
			Addr:    addr,
		},
		notify:          make(chan error, 1),
		shutdownTimeout: 5 * time.Second,
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
	case s := <-interrupt:
		fmt.Printf("Server run interrupted by OS signal %s", s.String())
	case err := <-s.Notify():
		fmt.Printf("Server connection error %v", err)
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
