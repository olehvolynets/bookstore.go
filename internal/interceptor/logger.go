package interceptor

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/felixge/httpsnoop"
	"github.com/rs/zerolog"
	"github.com/uptrace/bunrouter"

	"bookstore/internal/logging"
)

func ReqLogging(ctx context.Context) bunrouter.MiddlewareFunc {
	log := logging.WithLogSource(logging.FromContext(ctx), "interceptor")
	reqLog := logging.WithLogSource(logging.FromContext(ctx), "bookstore")

	log.Trace().Msg("enabling request logging")

	var reqFormatter logFn
	if os.Getenv("LOG_MODE") == "structured" {
		reqFormatter = structuredFormatter
	} else {
		reqFormatter = prettyFormatter
	}

	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, req bunrouter.Request) error {
			res := NewResponseWriter(w)

			now := time.Now()
			err := next(w, req)
			dur := time.Since(now)
			statusCode := res.StatusCode()

			var evt *zerolog.Event
			if statusCode < 500 {
				evt = reqLog.Info()
			} else {
				evt = reqLog.Error()
			}

			if err != nil {
				evt.Err(err)
			}

			reqFormatter(evt, statusCode, dur, req.Method, req.URL.String())

			return err
		}
	}
}

type logFn func(evt *zerolog.Event, status int, dur time.Duration, met string, url string)

func structuredFormatter(
	evt *zerolog.Event,
	status int,
	dur time.Duration,
	meth string,
	url string,
) {
	evt.Int("status", status).Dur("duration", dur).Str("method", meth).Str("url", url).Send()
}

func prettyFormatter(evt *zerolog.Event, status int, dur time.Duration, meth string, url string) {
	evt.Msg(
		fmt.Sprint(
			formatStatus(status),
			fmt.Sprintf(" %10s ", dur.Round(time.Microsecond)),
			formatMethod(meth),
			" ",
			url,
		),
	)
}

//------------------------------------------------------------------------------

type ResponseWriter struct {
	Wrapped    http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	var rw ResponseWriter
	rw.Wrapped = httpsnoop.Wrap(w, httpsnoop.Hooks{
		WriteHeader: func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return func(statusCode int) {
				if rw.statusCode == 0 {
					rw.statusCode = statusCode
				}
				next(statusCode)
			}
		},
	})
	return &rw
}

func (w *ResponseWriter) StatusCode() int {
	if w.statusCode != 0 {
		return w.statusCode
	}
	return http.StatusOK
}

//------------------------------------------------------------------------------

func formatStatus(code int) string {
	var col *color.Color

	switch {
	case code >= 200 && code < 300:
		col = color.New(color.BgGreen, color.FgHiWhite)
	case code >= 300 && code < 400:
		col = color.New(color.BgWhite, color.FgHiBlack)
	case code >= 400 && code < 500:
		col = color.New(color.BgYellow, color.FgHiBlack)
	default:
		col = color.New(color.BgRed, color.FgHiWhite)
	}

	return col.Sprintf(" %d ", code)
}

func formatMethod(method string) string {
	var col *color.Color

	switch method {
	case http.MethodGet:
		col = color.New(color.BgBlue, color.FgHiWhite)
	case http.MethodPost:
		col = color.New(color.BgGreen, color.FgHiWhite)
	case http.MethodPut:
		col = color.New(color.BgYellow, color.FgHiBlack)
	case http.MethodDelete:
		col = color.New(color.BgRed, color.FgHiWhite)
	case http.MethodPatch:
		col = color.New(color.BgCyan, color.FgHiWhite)
	case http.MethodHead:
		col = color.New(color.BgMagenta, color.FgHiWhite)
	default:
		col = color.New(color.BgWhite, color.FgHiBlack)
	}

	return col.Sprintf(" %-7s ", method)
}
