package internalhttp

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func simpleHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("ping-pong"))
}

func TestMiddleware(t *testing.T) {
	t.Run("loggingMiddleware", func(t *testing.T) {
		buf := &bytes.Buffer{}
		slog.SetDefault(slog.New(slog.NewTextHandler(buf, nil)))

		s := httptest.NewServer(loggingMiddleware(http.HandlerFunc(simpleHandler)))
		defer s.Close()

		res, err := http.Get(s.URL) //nolint:noctx
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusCreated)

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, body, []byte("ping-pong"))

		pattern := `time=[\d\-T:.Z]+ level=INFO msg=loggingMiddleware addr=[\d.:]+ method=GET path=\/ protocol=HTTP\/1.1 code=\d+ duration=\d+.+ user-agent=Go-http-client\/1.1` //nolint:lll
		require.Regexp(t, pattern, buf.String())
	})
}
