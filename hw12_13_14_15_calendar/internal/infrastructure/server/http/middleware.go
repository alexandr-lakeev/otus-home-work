package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
	appcalendar "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/google/uuid"
)

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func newLoggingMiddleware(logger app.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := uuid.New()

			ctx := context.WithValue(r.Context(), appcalendar.RequestIDContextKey, requestID)
			rw := &responseWriter{ResponseWriter: w}

			next.ServeHTTP(rw, r.Clone(ctx))

			logger.Info(fmt.Sprintf(
				"%s [%s] %s %s %s %d %v %s",
				r.RemoteAddr,
				start.Format(time.RFC3339),
				r.Method,
				r.URL,
				r.Proto,
				rw.code,
				time.Since(start),
				r.UserAgent(),
			))
		})
	}
}
