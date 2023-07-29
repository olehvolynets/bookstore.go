package server

import (
	"fmt"
	"net/http"

	"github.com/uptrace/bunrouter"

	"bookstore/log"
)

type Server struct {
	rootRouter *bunrouter.Router
	host       string
	port       uint
	store      Storage
}

func New(rr *bunrouter.Router, store Storage) *Server {
	s := &Server{
		rootRouter: rr,
		host:       "localhost",
		port:       3000,
		store:      store,
	}

	return s
}

func (s *Server) Start() error {
	log.Info().
		Str("host", s.host).
		Uint("port", s.port).
		Msg("listening")

	httpServer := http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
		Handler: s.rootRouter,
	}

	return httpServer.ListenAndServe()
}
