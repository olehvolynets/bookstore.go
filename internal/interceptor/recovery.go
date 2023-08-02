package interceptor

import (
	"net/http"

	"github.com/uptrace/bunrouter"

	"bookstore/internal/log"
)

func Recovery() bunrouter.MiddlewareFunc {
	log.Trace().Msg("enabling http handler recovery interceptor")

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
