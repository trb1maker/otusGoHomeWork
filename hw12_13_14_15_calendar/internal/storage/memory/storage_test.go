package memorystorage_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage/memory"
)

func TestStorage(t *testing.T) {
	store, err := memorystorage.New()
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var ownerID string

	t.Run("подключение", func(t *testing.T) {
		require.NoError(t, store.Connect(context.Background()))
	})

	t.Run("регистрация пользователя", func(t *testing.T) {
		ownerID, err = store.RegisterUser(ctx)
		require.NoError(t, err)
		require.NotZero(t, ownerID)
	})

	t.Run("нормальные сценарии", func(t *testing.T) {
		testEvent := storage.Event{
			Title:     "Test event",
			StartTime: time.Now().UTC(),
			EndTime:   time.Now().UTC().Add(time.Hour),
			OwnerID:   "62e1352f-f269-43ae-b784-272d9d8e4f8a",
			Notify:    time.Minute * 15,
		}

		eventID, err := store.InsertOne(ctx, testEvent)
		require.NoError(t, err)
		require.NotZero(t, eventID)

		testEvent.ID = eventID

		selectedEvent, err := store.SelectOne(ctx, eventID)
		require.NoError(t, err)
		require.Equal(t, testEvent, selectedEvent)

		testEvent.Description = "update"
		selectedEvent.Description = "update"
		require.NoError(t, store.UpdateOne(ctx, selectedEvent))

		selectedEvent, err = store.SelectOne(ctx, eventID)
		require.NoError(t, err)
		require.Equal(t, testEvent, selectedEvent)

		require.NoError(t, store.DeleteOne(ctx, eventID))

		selectedEvent, err = store.SelectOne(ctx, eventID)
		require.ErrorIs(t, err, storage.ErrNotFound)
		require.Zero(t, selectedEvent)
	})

	t.Run("негативные сценарии", func(t *testing.T) {
		require.ErrorIs(t, store.UpdateOne(ctx, storage.Event{ID: "abc"}), storage.ErrNotFound)
		require.ErrorIs(t, store.DeleteOne(ctx, "abc"), storage.ErrNotFound)
	})

	t.Run("сложные методы", func(t *testing.T) {
		events, err := store.SelectAllEvents(ctx, ownerID)
		require.NoError(t, err)
		require.Len(t, events, 0)

		for days := 1; days <= 10; days++ {
			event := storage.Event{
				Title:     fmt.Sprintf("Event %d", days),
				StartTime: time.Now().Add(time.Duration(-days) * 24 * time.Hour),
				EndTime:   time.Now().Add(time.Duration(-days) * 24 * time.Hour).Add(time.Hour),
				OwnerID:   ownerID,
			}

			_, err = store.InsertOne(ctx, event)
			require.NoError(t, err)
		}

		events, err = store.SelectAllEvents(ctx, ownerID)
		require.NoError(t, err)
		require.Len(t, events, 10)

		next, err := store.InsertOne(ctx, storage.Event{
			Title:     "NextEvent",
			StartTime: time.Now().Add(24 * time.Hour),
			EndTime:   time.Now().Add(25 * time.Hour),
			OwnerID:   ownerID,
		})
		require.NoError(t, err)

		events, err = store.SelectEventsBetweenDates(
			ctx,
			ownerID,
			time.Now().Add(-30*time.Hour),
			time.Now().Add(30*time.Hour),
		)
		require.NoError(t, err)
		require.Len(t, events, 2)

		nextEvent, _ := store.SelectNextEvent(ctx, ownerID)
		// require.NoError(t, err)
		require.Equal(t, next, nextEvent.ID)
	})

	t.Run("отключение", func(t *testing.T) {
		require.NoError(t, store.Close(ctx))
	})
}
