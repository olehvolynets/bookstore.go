package interceptor

import (
	"context"
	"net/http"

	"github.com/uptrace/bunrouter"

	"bookstore/internal/logging"
)

func Recovery(ctx context.Context) bunrouter.MiddlewareFunc {
	log := logging.WithLogSource(logging.FromContext(ctx), "interceptor")

	log.Trace().Msg("enabling http handler recovery")

	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, req bunrouter.Request) error {
			defer func() {
				if p := recover(); p != nil {
					log.Error().Interface("panic_source", p).Msg("http handler panic")
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			return next(w, req)
		}
	}
}
