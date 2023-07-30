package server

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/uptrace/bunrouter"

	"bookstore/api/rest"
	"bookstore/db"
	"bookstore/internal/interceptor"
	"bookstore/log"
	"bookstore/server"
)

var (
	host string
	port uint
)

func init() {
	Command.PersistentFlags().
		StringVarP(&host, "host", "H", "localhost", "Host the server is running on.")
	Command.PersistentFlags().UintVarP(&port, "port", "p", 3000, "Host the server is running on.")
}

var Command = &cobra.Command{
	Use:     "server",
	Aliases: []string{"s"},
	Short:   "Server management",
	Run:     commandCb,
}

func commandCb(_ *cobra.Command, _ []string) {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		done()
		if r := recover(); r != nil {
			log.Fatal().Interface("panic", r).Msg("application panic")
		}
	}()

	router := bunrouter.New(
		bunrouter.Use(interceptor.Recovery()),
		bunrouter.Use(interceptor.ReqLogging()),
	).Compat()

	router.WithGroup("/api", rest.Routes)

	dbConfig := db.NewConfig("postgres")
	db, err := db.New(ctx, dbConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create a db connection pool")
	}
	defer db.Close(ctx)

	s, err := server.New(router, server.WithPort(port))
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	err = s.Start(ctx)
	done()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Info().Msg("shutdown success")
}
