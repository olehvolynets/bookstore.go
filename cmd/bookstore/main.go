package main

import (
	"bookstore/internal/bookstore"
	"bookstore/internal/db"
	"bookstore/internal/interceptor"
	"bookstore/internal/logging"
	"bookstore/internal/server"
	"context"
	"flag"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/uptrace/bunrouter"
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
	appConfig, err := bookstore.NewConfig(ctx)
	if err != nil {
		appLog.Fatal().Err(err).Send()
	}
	// 2. Create env (stores db pointers and stuff)
	// 3. Create application server(s) (for books, authros, etc.), it will generate http.Handler
	// 4. Create http server and pass http.Handler from step 3

	router := bunrouter.New(
		bunrouter.Use(interceptor.Recovery(ctx)),
		bunrouter.Use(interceptor.ReqLogging(ctx)),
	).Compat()

	wd, _ := os.Getwd()
	fs := http.FileServer(http.Dir(wd + "/web/static"))
	router.GET("/static/*file", http.StripPrefix("/static/", fs).ServeHTTP)

	router.GET("/", func(w http.ResponseWriter, r *http.Request) {
		temp := template.Must(template.ParseFiles(wd + "/cmd/bookstore/index.html"))
		err := temp.Execute(w, struct {
			Foo uint
			Els []string
			F   any
		}{
			Foo: 123,
			Els: []string{"asd", "foo", "bar"},
			F:   true,
		})
		if err != nil {
			appLog.Warn().Err(err).Send()
		}
	})

	db, err := db.New(ctx, &appConfig.Database)
	if err != nil {
		appLog.Fatal().Err(err).Msg("unable to create a db connection pool")
	}
	defer db.Close(ctx)

	tcpServer, err := server.New(server.WithPort(port))
	if err != nil {
		appLog.Fatal().Err(err).Send()
	}

	err = tcpServer.Start(ctx, router)
	done()
	if err != nil {
		appLog.Fatal().Err(err).Send()
	}

	appLog.Info().Msg("shutdown success")
}
