package internalgrpc

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/server/grpc/api"
	memorystorage "github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestServer(t *testing.T) {
	buf := &bytes.Buffer{}
	slog.SetDefault(slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	store, err := memorystorage.New()
	require.NoError(t, err)
	defer store.Close(ctx)

	require.NoError(t, store.Connect(ctx))

	userID, err := store.RegisterUser(ctx)
	require.NoError(t, err)

	server := NewServer(app.New(store), "localhost", 8081)

	go func() {
		require.NoError(t, server.Start(ctx))
	}()

	go func() {
		conn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)
		defer conn.Close()

		event := &api.Event{
			Title:     "created event",
			StartTime: timestamppb.New(time.Now()),
			EndTime:   timestamppb.New(time.Now().Add(30 * time.Minute)),
			Owner:     userID,
		}

		client := api.NewEventServiceClient(conn)
		resp, err := client.NewEvent(ctx, event)
		require.NoError(t, err)

		event.Id = resp.EventId
		require.NotZero(t, event.GetId())

		event.Description = "new description"

		_, err = client.UpdateEvent(ctx, event)
		require.NoError(t, err)

		events, err := client.All(ctx, &api.Request{UserId: userID})
		require.NoError(t, err)
		require.Len(t, events.Events, 1)

		_, err = client.DeleteEvent(ctx, &api.Request{EventId: event.Id})
		require.NoError(t, err)

		cancel()
	}()

	<-ctx.Done()

	logResult := buf.String()
	require.NotZero(t, logResult)
}
