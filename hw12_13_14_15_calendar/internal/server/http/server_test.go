package internalhttp_test

import (
	"context"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	internalhttp "github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/server/http"
)

func TestServer(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	s := internalhttp.NewServer(nil, "localhost", 12_000)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		defer wg.Done()

		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		require.NoError(t, s.Stop(ctx))
	}()

	go func() {
		defer wg.Done()

		//nolint:noctx
		res, err := http.Get("http://localhost:12000/ping")
		require.NoError(t, err)
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, res.StatusCode, http.StatusOK)
		require.Equal(t, body, []byte("pong"))
	}()

	require.NoError(t, s.Start(ctx))
	wg.Wait()
}
