package internalhttp

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/require"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage/memory"
)

func TestPing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	res := httptest.NewRecorder()

	s := &Server{}
	s.ping(res, req)

	require.Equal(t, res.Result().StatusCode, http.StatusOK) //nolint: bodyclose
	require.Equal(t, res.Body.String(), "pong")
}

func TestInternalHandlers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store, err := memorystorage.New()
	require.NoError(t, err)
	require.NoError(t, store.Connect(ctx))

	userID, err := store.RegisterUser(ctx)
	require.NoError(t, err)

	srv := NewServer(app.New(store), "localhost", 12000)

	startEvent := storage.Event{
		Title:       "start event",
		StartTime:   time.Now().Add(30 * time.Minute).UTC(),
		EndTime:     time.Now().Add(60 * time.Minute).UTC(),
		Description: "test event",
		OwnerID:     userID,
		Notify:      5 * time.Minute,
	}
	buf := &bytes.Buffer{}

	t.Run("create event", func(t *testing.T) {
		_, err := easyjson.MarshalToWriter(startEvent, buf)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/event", buf)
		w := httptest.NewRecorder()

		srv.postEvent(w, req)
		require.Equal(t, w.Result().StatusCode, http.StatusOK) //nolint: bodyclose

		var obj dto
		require.NoError(t, easyjson.UnmarshalFromReader(w.Body, &obj))
		require.True(t, obj.Ok)
		require.Equal(t, obj.Count, 1)
		require.Len(t, obj.Events, 1)

		startEvent.ID = obj.Events[0].ID

		buf.Reset()
	})

	t.Run("get event", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/event/{eventID}", nil)
		req.SetPathValue("eventID", startEvent.ID)
		w := httptest.NewRecorder()

		srv.getEvent(w, req)
		require.Equal(t, w.Result().StatusCode, http.StatusOK) //nolint: bodyclose

		var obj dto
		require.NoError(t, easyjson.UnmarshalFromReader(w.Body, &obj))
		require.True(t, obj.Ok)
		require.Equal(t, obj.Count, 1)
		require.Len(t, obj.Events, 1)
	})

	t.Run("update event", func(t *testing.T) {
		startEvent.Notify = 10 * time.Minute
		_, err := easyjson.MarshalToWriter(startEvent, buf)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPut, "/event", buf)
		w := httptest.NewRecorder()

		srv.putEvent(w, req)
		require.Equal(t, w.Result().StatusCode, http.StatusOK) //nolint: bodyclose
		var obj dto
		require.NoError(t, easyjson.UnmarshalFromReader(w.Body, &obj))
		require.True(t, obj.Ok)

		buf.Reset()
	})

	t.Run("delete event", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/event/{eventID}", nil)
		req.SetPathValue("eventID", startEvent.ID)
		w := httptest.NewRecorder()

		srv.deleteEvent(w, req)
		require.Equal(t, w.Result().StatusCode, http.StatusOK) //nolint: bodyclose

		var obj dto
		require.NoError(t, easyjson.UnmarshalFromReader(w.Body, &obj))
		require.True(t, obj.Ok)
	})

	t.Run("get all events", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			easyjson.MarshalToWriter(startEvent, buf)

			req := httptest.NewRequest(http.MethodPost, "/event", buf)
			w := httptest.NewRecorder()

			srv.postEvent(w, req)
			require.Equal(t, w.Result().StatusCode, http.StatusOK) //nolint: bodyclose

			buf.Reset()
		}

		req := httptest.NewRequest(http.MethodGet, "/user/{userID}/all", nil)
		req.SetPathValue("userID", userID)
		w := httptest.NewRecorder()

		srv.getAllEvents(w, req)
		require.Equal(t, w.Result().StatusCode, http.StatusOK) //nolint: bodyclose

		var obj dto
		require.NoError(t, easyjson.UnmarshalFromReader(w.Body, &obj))
		require.True(t, obj.Ok)
		require.Equal(t, obj.Count, 5)
		require.Len(t, obj.Events, 5)

		buf.Reset()
	})

	t.Run("get next event", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/{userID}/next", nil)
		req.SetPathValue("userID", userID)
		w := httptest.NewRecorder()

		srv.getNextEvent(w, req)
		require.Equal(t, w.Result().StatusCode, http.StatusOK) //nolint: bodyclose

		var obj dto
		require.NoError(t, easyjson.UnmarshalFromReader(w.Body, &obj))
		require.True(t, obj.Ok)
		require.Equal(t, obj.Count, 1)
		require.Len(t, obj.Events, 1)
	})

	t.Run("get current day", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/{userID}/day", nil)
		req.SetPathValue("userID", userID)
		w := httptest.NewRecorder()

		srv.getDayEvents(w, req)
		require.Equal(t, w.Result().StatusCode, http.StatusOK) //nolint: bodyclose

		var obj dto
		require.NoError(t, easyjson.UnmarshalFromReader(w.Body, &obj))
		require.True(t, obj.Ok)
		require.Equal(t, obj.Count, 5)
		require.Len(t, obj.Events, 5)
	})

	t.Run("get current week", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/{userID}/week", nil)
		req.SetPathValue("userID", userID)
		w := httptest.NewRecorder()

		srv.getWeekEvents(w, req)
		require.Equal(t, w.Result().StatusCode, http.StatusOK) //nolint: bodyclose

		var obj dto
		require.NoError(t, easyjson.UnmarshalFromReader(w.Body, &obj))
		require.True(t, obj.Ok)
		require.Equal(t, obj.Count, 5)
		require.Len(t, obj.Events, 5)
	})

	t.Run("get current month", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/{userID}/month", nil)
		req.SetPathValue("userID", userID)
		w := httptest.NewRecorder()

		srv.getMonthEvents(w, req)
		require.Equal(t, w.Result().StatusCode, http.StatusOK) //nolint: bodyclose

		var obj dto
		require.NoError(t, easyjson.UnmarshalFromReader(w.Body, &obj))
		require.True(t, obj.Ok)
		require.Equal(t, obj.Count, 5)
		require.Len(t, obj.Events, 5)
	})

	store.Close(context.Background())
}

func TestUserMethods(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store, err := memorystorage.New()
	require.NoError(t, err)
	require.NoError(t, store.Connect(ctx))

	userID, err := store.RegisterUser(ctx)
	require.NoError(t, err)

	srv := NewServer(app.New(store), "localhost", 12000)
	buf := &bytes.Buffer{}

	events := []storage.Event{
		{
			Title:     "test event",
			StartTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC),
			OwnerID:   userID,
		},
		{
			Title:     "test event",
			StartTime: time.Date(2023, 1, 1, 15, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2023, 1, 1, 16, 0, 0, 0, time.UTC),
			OwnerID:   userID,
		},
		{
			Title:     "test event",
			StartTime: time.Date(2023, 1, 14, 15, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2023, 1, 14, 16, 0, 0, 0, time.UTC),
			OwnerID:   userID,
		},
	}

	for i := 0; i < len(events); i++ {
		_, err := easyjson.MarshalToWriter(events[i], buf)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/event", buf)
		w := httptest.NewRecorder()

		srv.postEvent(w, req)
		require.Equal(t, w.Result().StatusCode, http.StatusOK) //nolint: bodyclose
		buf.Reset()
	}

	// События за день.
	req := httptest.NewRequest(http.MethodGet, "/user/{userID}/day?start=2023-01-01", nil)
	req.SetPathValue("userID", userID)
	w := httptest.NewRecorder()

	srv.getDayEvents(w, req)
	require.Equal(t, w.Result().StatusCode, http.StatusOK) //nolint: bodyclose

	var obj dto
	require.NoError(t, easyjson.UnmarshalFromReader(w.Body, &obj))
	require.True(t, obj.Ok)
	require.Equal(t, obj.Count, 2)
	require.Len(t, obj.Events, 2)

	// События за неделю
	req = httptest.NewRequest(http.MethodGet, "/user/{userID}/week?start=2022-12-26", nil)
	req.SetPathValue("userID", userID)
	w = httptest.NewRecorder()

	srv.getWeekEvents(w, req)
	require.Equal(t, w.Result().StatusCode, http.StatusOK) //nolint: bodyclose

	require.NoError(t, easyjson.UnmarshalFromReader(w.Body, &obj))
	require.True(t, obj.Ok)
	require.Equal(t, obj.Count, 2)
	require.Len(t, obj.Events, 2)

	// События за месяц
	req = httptest.NewRequest(http.MethodGet, "/user/{userID}/week?start=2023-01-01", nil)
	req.SetPathValue("userID", userID)
	w = httptest.NewRecorder()

	srv.getMonthEvents(w, req)
	require.Equal(t, w.Result().StatusCode, http.StatusOK) //nolint: bodyclose

	require.NoError(t, easyjson.UnmarshalFromReader(w.Body, &obj))
	require.True(t, obj.Ok)
	require.Equal(t, obj.Count, 3)
	require.Len(t, obj.Events, 3)

	store.Close(context.Background())
}
