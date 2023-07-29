package main

import (
	"github.com/uptrace/bunrouter"

	"bookstore/api/rest"
	"bookstore/log"
	"bookstore/server"
)

func main() {
	router := bunrouter.New(
		bunrouter.Use(log.NewMiddleware()),
	)

	router.WithGroup("/api", rest.Routes)

	s := server.New(router, nil)

	log.Fatal().Err(s.Start()).Send()
}
