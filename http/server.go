package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func NewServer(handler http.Handler, addr string) (*Server, error) {
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
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
