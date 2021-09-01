package logging

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Middleware(next http.Handler) http.Handler {
	middleware := func(w http.ResponseWriter, r *http.Request) {
		printRequest(r)

		lrw := &loggingResponseWriter{w, 0}

		startTime := time.Now()
		next.ServeHTTP(lrw, r)

		if lrw.status == 0 {
			lrw.WriteHeader(http.StatusOK)
		}

		printResponse(lrw, startTime)
	}

	return http.HandlerFunc(middleware)
}

func printRequest(r *http.Request) {
	log.Info().
		Str("Address", r.RemoteAddr).
		Str("Method", r.Method).
		Str("Url", r.URL.String()).
		Msg("Serving route")
}

func printResponse(w *loggingResponseWriter, start time.Time) {
	log.Info().
		Int("Response", w.status).
		TimeDiff("Duration", time.Now(), start).
		Msg("Route served")
}
