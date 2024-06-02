package sqlstorage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

// TODO: Сделать тестовое окружение, поднимаемое автоматически в docker.
func TestStorage(t *testing.T) {
	t.SkipNow() // Проверяю только в готовом окружении

	store, err := New("192.168.0.103", 5432, "otus", "otus", "20240518_otus")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var ownerID string

	t.Run("подключение", func(t *testing.T) {
		require.NoError(t, store.Connect(ctx))
	})

	t.Run("регистрация нового пользователя", func(t *testing.T) {
		ownerID, err = store.RegisterUser(ctx)
		require.NoError(t, err)
	})

	t.Run("нормальные сценарии", func(t *testing.T) {
		var (
			startTime = time.Now().UTC()
			endTime   = startTime.Add(time.Hour)
			notify    = endTime.Sub(startTime)
		)

		e := storage.Event{
			Title:     "Event",
			StartTime: startTime,
			EndTime:   endTime,
			OwnerID:   ownerID,
			Notify:    notify,
		}

		eventID, err := store.InsertOne(ctx, e)
		require.NoError(t, err)
		require.NotZero(t, eventID)

		e.ID = eventID
		e1, err := store.SelectOne(ctx, eventID)
		require.NoError(t, err)
		require.Equal(t, e, e1)

		e.Description = "update"
		require.NoError(t, store.UpdateOne(ctx, e))

		require.NoError(t, store.DeleteOne(ctx, eventID))
	})

	t.Run("событие не найдено", func(t *testing.T) {
		_, err := store.SelectOne(ctx, "62e1352f-f269-43ae-b784-272d9d8e4f8a")
		require.ErrorIs(t, err, storage.ErrNotFound)

		require.ErrorIs(t,
			store.UpdateOne(ctx, storage.Event{ID: "62e1352f-f269-43ae-b784-272d9d8e4f8a"}),
			storage.ErrNotFound,
		)

		require.ErrorIs(t, store.DeleteOne(ctx, "62e1352f-f269-43ae-b784-272d9d8e4f8a"), storage.ErrNotFound)
	})

	t.Run("отключение", func(t *testing.T) {
		require.NoError(t, store.Close(ctx))
	})
}
