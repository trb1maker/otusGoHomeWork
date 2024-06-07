package internalhttp

import (
	"errors"
	"net/http"
	"time"

	"github.com/mailru/easyjson"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

func (s *Server) ping(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("pong"))
}

func (s *Server) postEvent(w http.ResponseWriter, r *http.Request) {
	var event storage.Event

	if err := easyjson.UnmarshalFromReader(r.Body, &event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		easyjson.MarshalToHTTPResponseWriter(dto{Details: err.Error()}, w)
		return
	}

	eventID, err := s.app.CreateEvent(
		r.Context(),
		event.OwnerID,
		event.Title,
		event.StartTime.UTC(),
		event.EndTime.UTC(),
		event.Description,
		event.Notify,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		easyjson.MarshalToHTTPResponseWriter(dto{Details: err.Error()}, w)
		return
	}

	event.ID = eventID
	obj := dto{
		Ok:     true,
		Count:  1,
		Events: []storage.Event{event},
	}

	easyjson.MarshalToHTTPResponseWriter(obj, w)
}

func (s *Server) getEvent(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("eventID")

	event, err := s.app.GetEvent(r.Context(), eventID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			easyjson.MarshalToHTTPResponseWriter(dto{Details: err.Error()}, w)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		easyjson.MarshalToHTTPResponseWriter(dto{Details: err.Error()}, w)
		return
	}

	obj := dto{
		Ok:     true,
		Count:  1,
		Events: []storage.Event{event},
	}
	easyjson.MarshalToHTTPResponseWriter(obj, w)
}

func (s *Server) putEvent(w http.ResponseWriter, r *http.Request) {
	var event storage.Event

	if err := easyjson.UnmarshalFromReader(r.Body, &event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		easyjson.MarshalToHTTPResponseWriter(dto{Details: err.Error()}, w)
		return
	}

	if err := s.app.UpdateEvent(r.Context(), event); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			easyjson.MarshalToHTTPResponseWriter(dto{Details: err.Error()}, w)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		easyjson.MarshalToHTTPResponseWriter(dto{Details: err.Error()}, w)
		return
	}

	easyjson.MarshalToHTTPResponseWriter(dto{Ok: true}, w)
}

func (s *Server) deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("eventID")

	if err := s.app.DeleteEvent(r.Context(), eventID); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			easyjson.MarshalToHTTPResponseWriter(dto{Details: err.Error()}, w)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		easyjson.MarshalToHTTPResponseWriter(dto{Details: err.Error()}, w)
		return
	}

	easyjson.MarshalToHTTPResponseWriter(dto{Ok: true}, w)
}

func (s *Server) getAllEvents(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")

	events, err := s.app.GetAllEvents(r.Context(), userID)
	if err != nil {
		if errors.Is(err, storage.ErrNoData) {
			easyjson.MarshalToHTTPResponseWriter(dto{Ok: true}, w)
			return
		}
		easyjson.MarshalToHTTPResponseWriter(dto{Ok: false, Details: err.Error()}, w)
		return
	}
	easyjson.MarshalToHTTPResponseWriter(dto{Ok: true, Count: len(events), Events: events}, w)
}

func (s *Server) getNextEvent(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")

	event, err := s.app.GetNextEvent(r.Context(), userID)
	if err != nil {
		if errors.Is(err, storage.ErrNoData) {
			easyjson.MarshalToHTTPResponseWriter(dto{Ok: true}, w)
			return
		}
		easyjson.MarshalToHTTPResponseWriter(dto{Ok: false, Details: err.Error()}, w)
		return
	}
	easyjson.MarshalToHTTPResponseWriter(dto{Ok: true, Count: 1, Events: []storage.Event{event}}, w)
}

//nolint:all
func (s *Server) getDayEvents(w http.ResponseWriter, r *http.Request) {
	var (
		events []storage.Event
		err    error
	)

	userID := r.PathValue("userID")

	startOption := r.FormValue("start")
	if startOption == "" {
		events, err = s.app.GetEventsCurrentDay(r.Context(), userID)
	} else {
		startTime, err := time.Parse(time.DateOnly, startOption)
		if err != nil {
			easyjson.MarshalToHTTPResponseWriter(dto{Ok: false, Details: err.Error()}, w)
			return
		}
		events, err = s.app.GetEventsDayAfter(r.Context(), userID, startTime)
	}
	if err != nil {
		if errors.Is(err, storage.ErrNoData) {
			easyjson.MarshalToHTTPResponseWriter(dto{Ok: true}, w)
			return
		}
		easyjson.MarshalToHTTPResponseWriter(dto{Ok: false, Details: err.Error()}, w)
		return
	}

	easyjson.MarshalToHTTPResponseWriter(dto{Ok: true, Count: len(events), Events: events}, w)
}

//nolint:all
func (s *Server) getWeekEvents(w http.ResponseWriter, r *http.Request) {
	var (
		events []storage.Event
		err    error
	)

	userID := r.PathValue("userID")

	startOption := r.FormValue("start")
	if startOption == "" {
		events, err = s.app.GetEventsCurrentWeek(r.Context(), userID)
	} else {
		startTime, err := time.Parse(time.DateOnly, startOption)
		if err != nil {
			easyjson.MarshalToHTTPResponseWriter(dto{Ok: false, Details: err.Error()}, w)
			return
		}
		events, err = s.app.GetEventsWeekAfter(r.Context(), userID, startTime)
	}
	if err != nil {
		if errors.Is(err, storage.ErrNoData) {
			easyjson.MarshalToHTTPResponseWriter(dto{Ok: true}, w)
			return
		}
		easyjson.MarshalToHTTPResponseWriter(dto{Ok: false, Details: err.Error()}, w)
		return
	}

	easyjson.MarshalToHTTPResponseWriter(dto{Ok: true, Count: len(events), Events: events}, w)
}

//nolint:all
func (s *Server) getMonthEvents(w http.ResponseWriter, r *http.Request) {
	var (
		events []storage.Event
		err    error
	)

	userID := r.PathValue("userID")

	startOption := r.FormValue("start")
	if startOption == "" {
		events, err = s.app.GetEventsCurrentMonth(r.Context(), userID)
	} else {
		startTime, err := time.Parse(time.DateOnly, startOption)
		if err != nil {
			easyjson.MarshalToHTTPResponseWriter(dto{Ok: false, Details: err.Error()}, w)
			return
		}
		events, err = s.app.GetEventsMonthAfter(r.Context(), userID, startTime)
	}
	if err != nil {
		if errors.Is(err, storage.ErrNoData) {
			easyjson.MarshalToHTTPResponseWriter(dto{Ok: true}, w)
			return
		}
		easyjson.MarshalToHTTPResponseWriter(dto{Ok: false, Details: err.Error()}, w)
		return
	}

	easyjson.MarshalToHTTPResponseWriter(dto{Ok: true, Count: len(events), Events: events}, w)
}
