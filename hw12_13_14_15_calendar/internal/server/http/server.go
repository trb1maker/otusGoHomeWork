package internalhttp

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	app Application
	srv *http.Server
}

type Application interface {
	CreateEvent(ctx context.Context, userID string, title string, startTime time.Time,
		endTime time.Time, description string, notify time.Duration) (string, error)
	GetEvent(ctx context.Context, eventID string) (storage.Event, error)
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, eventID string) error
	GetNextEvent(ctx context.Context, ownerID string) (storage.Event, error)
	GetAllEvents(ctx context.Context, ownerID string) ([]storage.Event, error)
	GetEventsFromRange(ctx context.Context, ownerID string, from time.Time, to time.Time) ([]storage.Event, error)
	GetEventsCurrentDay(ctx context.Context, ownerID string) ([]storage.Event, error)
	GetEventsCurrentWeek(ctx context.Context, ownerID string) ([]storage.Event, error)
	GetEventsCurrentMonth(ctx context.Context, ownerID string) ([]storage.Event, error)
	GetEventsDayAfter(ctx context.Context, ownerID string, start time.Time) ([]storage.Event, error)
	GetEventsWeekAfter(ctx context.Context, ownerID string, start time.Time) ([]storage.Event, error)
	GetEventsMonthAfter(ctx context.Context, ownerID string, start time.Time) ([]storage.Event, error)
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

	mux.Handle("POST /event", loggingMiddleware(http.HandlerFunc(s.postEvent)))
	mux.Handle("GET /event/{eventID}", loggingMiddleware(http.HandlerFunc(s.getEvent)))
	mux.Handle("PUT /event", loggingMiddleware(http.HandlerFunc(s.putEvent)))
	mux.Handle("DELETE /event/{eventID}", loggingMiddleware(http.HandlerFunc(s.deleteEvent)))

	mux.Handle("GET /user/{userID}/all", loggingMiddleware(http.HandlerFunc(s.getAllEvents)))
	mux.Handle("GET /user/{userID}/next", loggingMiddleware(http.HandlerFunc(s.getNextEvent)))
	mux.Handle("GET /user/{userID}/day", loggingMiddleware(http.HandlerFunc(s.getDayEvents)))
	mux.Handle("GET /user/{userID}/week", loggingMiddleware(http.HandlerFunc(s.getWeekEvents)))
	mux.Handle("GET /user/{userID}/month", loggingMiddleware(http.HandlerFunc(s.getMonthEvents)))

	return s
}

func (s *Server) Start(ctx context.Context) error {
	s.srv.BaseContext = func(_ net.Listener) context.Context {
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
