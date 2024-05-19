package sqlstorage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

// TODO: Сделать тестовое окружение, поднимаемое автоматически в docker.
func TestStorage(t *testing.T) {
	t.SkipNow() // Проверяю только в готовом окружении

	config := map[string]string{
		"DBHOST":     "192.168.0.103",
		"DBPORT":     "5432",
		"DBNAME":     "otus",
		"DBUSER":     "otus",
		"DBPASSWORD": "20240518_otus",
	}

	for name, value := range config {
		require.NoError(t, os.Setenv(name, value))
	}

	s := New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("подключение", func(t *testing.T) {
		require.NoError(t, s.Connect(ctx))
	})

	t.Run("нормальные сценарии", func(t *testing.T) {
		var (
			ownerID   = "ef6acd9c-7385-420d-b408-f0029a9decc3"
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

		eventID, err := s.InsertOne(ctx, e)
		require.NoError(t, err)
		require.NotZero(t, eventID)

		e.ID = eventID
		e1, err := s.SelectOne(ctx, eventID)
		require.NoError(t, err)
		require.Equal(t, e, e1)

		e.Description = "update"
		require.NoError(t, s.UpdateOne(ctx, e))

		require.NoError(t, s.DeleteOne(ctx, eventID))
	})

	t.Run("событие не найдено", func(t *testing.T) {
		_, err := s.SelectOne(ctx, "62e1352f-f269-43ae-b784-272d9d8e4f8a")
		require.ErrorIs(t, err, storage.ErrNotFound)

		require.ErrorIs(t, s.UpdateOne(ctx, storage.Event{ID: "62e1352f-f269-43ae-b784-272d9d8e4f8a"}), storage.ErrNotFound)

		require.ErrorIs(t, s.DeleteOne(ctx, "62e1352f-f269-43ae-b784-272d9d8e4f8a"), storage.ErrNotFound)
	})

	t.Run("отключение", func(t *testing.T) {
		require.NoError(t, s.Close(ctx))
	})
}
