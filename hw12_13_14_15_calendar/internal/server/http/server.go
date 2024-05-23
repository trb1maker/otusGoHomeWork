package internalhttp

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	app Application
	srv *http.Server
}

type Application interface {
	// TODO
}

func NewServer(app Application, host string, port int) *Server {
	mux := http.NewServeMux()

	s := &Server{
		app: app,
		srv: &http.Server{
			Addr:              net.JoinHostPort(host, strconv.Itoa(port)),
			Handler:           mux,
			ReadHeaderTimeout: time.Second,
			ErrorLog:          slog.NewLogLogger(slog.Default().Handler(), slog.LevelError),
		},
	}

	mux.Handle("/ping", loggingMiddleware(http.HandlerFunc(s.ping)))

	return s
}

func (s *Server) Start(ctx context.Context) error {
	s.srv.BaseContext = func(l net.Listener) context.Context {
		return context.WithoutCancel(ctx)
	}

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
