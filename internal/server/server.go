package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"bookstore/internal/log"
)

type Server struct {
	ip       string
	port     uint
	listener net.Listener
}

type ServerCfgOption func(*Server)

func New(opts ...ServerCfgOption) (*Server, error) {
	s := &Server{port: 3000}

	for _, opt := range opts {
		opt(s)
	}

	addr := fmt.Sprintf(":%d", s.port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create a listener on %s: %w", addr, err)
	}

	s.ip = listener.Addr().(*net.TCPAddr).IP.String()
	s.listener = listener

	return s, nil
}

func (s *Server) Start(ctx context.Context, handler http.Handler) error {
	httpServer := http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           handler,
	}

	errCh := make(chan error, 1)

	log.Info().
		Str("ip", s.ip).
		Uint("port", s.port).
		Msg("listening")

	go func() {
		<-ctx.Done()

		log.Debug().Msg("server.Start context closed")
		shutdownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		log.Debug().Msg("server.Start shutting down")
		errCh <- httpServer.Shutdown(shutdownCtx)
	}()

	if err := httpServer.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func WithPort(port uint) ServerCfgOption {
	return func(s *Server) {
		s.port = port
	}
}
