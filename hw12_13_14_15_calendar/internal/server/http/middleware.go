package internalhttp

import (
	"net/http"
	"time"

	"golang.org/x/exp/slog"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusWriter) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		sw := &statusWriter{w, http.StatusInternalServerError}
		next.ServeHTTP(sw, r)

		slog.LogAttrs(
			r.Context(),
			slog.LevelInfo,
			"loggingMiddleware",
			slog.String("addr", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("protocol", r.Proto),
			slog.Int("code", sw.status),
			slog.String("duration", time.Since(start).String()),
			slog.String("user-agent", r.UserAgent()),
		)
	})
}
