package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

func TestIntegrationApp(t *testing.T) {
	store, err := New()
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	require.NoError(t, store.Connect(ctx))

	userID := "62e1352f-f269-43ae-b784-272d9d8e4f8a"
	_, err = store.db.ExecContext(ctx, "insert into users values ($1)", userID)
	require.NoError(t, err)

	app := app.New(store)

	t.Run("создание событий", func(t *testing.T) {
		events := []storage.Event{
			{
				Title:     "Current week 1",
				StartTime: time.Date(2024, 5, 24, 12, 34, 56, 0, time.UTC),
				EndTime:   time.Date(2024, 5, 24, 12, 34, 56, 0, time.UTC),
				OwnerID:   userID,
			},
			{
				Title:     "Current week 2",
				StartTime: time.Date(2024, 5, 21, 12, 34, 56, 0, time.UTC),
				EndTime:   time.Date(2024, 5, 21, 12, 34, 56, 0, time.UTC),
				OwnerID:   userID,
			},
			{
				Title:     "Current month 1",
				StartTime: time.Date(2024, 5, 21, 12, 34, 56, 0, time.UTC),
				EndTime:   time.Date(2024, 5, 21, 12, 34, 56, 0, time.UTC),
				OwnerID:   userID,
			},
			{
				Title:     "Current month 2",
				StartTime: time.Date(2024, 5, 14, 12, 34, 56, 0, time.UTC),
				EndTime:   time.Date(2024, 5, 14, 12, 34, 56, 0, time.UTC),
				OwnerID:   userID,
			},
			{
				Title:     "Current month 3",
				StartTime: time.Date(2024, 5, 29, 12, 34, 56, 0, time.UTC),
				EndTime:   time.Date(2024, 5, 29, 12, 34, 56, 0, time.UTC),
				OwnerID:   userID,
			},
			{
				Title:     "Current day",
				StartTime: time.Date(2024, 5, 23, 12, 34, 56, 0, time.UTC),
				EndTime:   time.Date(2024, 5, 23, 12, 34, 56, 0, time.UTC),
				OwnerID:   userID,
			},
		}

		for _, event := range events {
			_, err = app.CreateEvent(ctx, userID, event.Title, event.StartTime, event.EndTime, event.Description, event.Notify)
			require.NoError(t, err)
		}
	})

	t.Run("события за текущий день", func(t *testing.T) {
		events, err := app.GetEventsCurrentDay(ctx, userID)
		require.NoError(t, err)
		require.Len(t, events, 1)
	})

	t.Run("события за текущий неделю", func(t *testing.T) {
		events, err := app.GetEventsCurrentWeek(ctx, userID)
		require.NoError(t, err)
		require.Len(t, events, 4)
	})

	t.Run("события за текущий месяц", func(t *testing.T) {
		events, err := app.GetEventsCurrentMonth(ctx, userID)
		require.NoError(t, err)
		require.Len(t, events, 6)
	})

	require.NoError(t, store.Close(ctx))
}
