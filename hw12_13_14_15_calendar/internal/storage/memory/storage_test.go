package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

func TestStorage(t *testing.T) {
	s := &Storage{
		store: make(map[string]storage.Event),
	}

	var (
		ownerID   = "62e1352f-f269-43ae-b784-272d9d8e4f8a"
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

	ctx := context.TODO()

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

	e1, err = s.SelectOne(ctx, "abc")
	require.ErrorIs(t, err, storage.ErrNotFound)
	require.Zero(t, e1)

	require.ErrorIs(t, s.UpdateOne(ctx, storage.Event{ID: "abc"}), storage.ErrNotFound)
	require.ErrorIs(t, s.DeleteOne(ctx, "abc"), storage.ErrNotFound)
}
