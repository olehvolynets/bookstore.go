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
	"bookstore/internal/logging"
	"bookstore/internal/server"
)

var port uint = *flag.Uint("port", 3000, "port for the server")

func main() {
	flag.Parse()

	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	logger := logging.NewLoggerFromEnv()
	ctx = logging.CtxWithLogger(ctx, logger)

	appLog := logging.WithLogSource(logger, "bookstore")

	defer func() {
		done()
		if r := recover(); r != nil {
			appLog.Fatal().Interface("panic", r).Msg("application panic")
		}
	}()

	// 1. Create config (parse env vars and create a struct)
	// 2. Create env (stores db pointers and stuff)
	// 3. Create application server(s) (for books, authros, etc.), it will generate http.Handler
	// 4. Create http server and pass http.Handler from step 3

	router := bunrouter.New(
		bunrouter.Use(interceptor.Recovery(ctx)),
		bunrouter.Use(interceptor.ReqLogging(ctx)),
	).Compat()

	router.WithGroup("/api", rest.Routes)

	dbConfig := db.NewConfig(ctx, "postgres")
	db, err := db.New(ctx, dbConfig)
	if err != nil {
		appLog.Fatal().Err(err).Msg("unable to create a db connection pool")
	}
	defer db.Close(ctx)

	s, err := server.New(server.WithPort(port))
	if err != nil {
		appLog.Fatal().Err(err).Send()
	}

	err = s.Start(ctx, router)
	done()
	if err != nil {
		appLog.Fatal().Err(err).Send()
	}

	appLog.Info().Msg("shutdown success")
}
