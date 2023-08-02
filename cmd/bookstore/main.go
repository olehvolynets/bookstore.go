package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/uptrace/bunrouter"

	"bookstore/api/rest"
	"bookstore/internal/db"
	"bookstore/internal/interceptor"
	"bookstore/internal/log"
	"bookstore/internal/server"
)

var port uint = *flag.Uint("port", 3000, "port for the server")

func main() {
	flag.Parse()

	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		done()
		if r := recover(); r != nil {
			log.Fatal().Interface("panic", r).Msg("application panic")
		}
	}()

	// 1. Create config (parse env vars and create a struct)
	// 2. Create env (stores db pointers and stuff)
	// 3. Create application server(s) (for books, authros, etc.), it will generate http.Handler
	// 4. Create http server and pass http.Handler from step 3

	router := bunrouter.New(
		bunrouter.Use(interceptor.Recovery()),
		bunrouter.Use(interceptor.ReqLogging()),
	).Compat()

	router.WithGroup("/api", rest.Routes)

	dbConfig := db.NewConfig(ctx, "postgres")
	db, err := db.New(ctx, dbConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create a db connection pool")
	}
	defer db.Close(ctx)

	s, err := server.New(server.WithPort(port))
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	err = s.Start(ctx, router)
	done()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Info().Msg("shutdown success")
}
