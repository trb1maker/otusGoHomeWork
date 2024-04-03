package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_copyWithOffset(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		tests := []struct {
			name   string
			srs    string
			offset int64
			limit  int64
			want   string
		}{
			{"read from start", "1234567890", 0, 5, "12345"},
			{"read from middle", "1234567890", 3, 5, "45678"},
			{"read from end", "1234567890", 5, 5, "67890"},
			{"limit great then size", "1234567890", 0, 15, "1234567890"},
			{"limit great then size with offset", "1234567890", 3, 15, "4567890"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var (
					r = bytes.NewReader([]byte(tt.srs))
					w = bytes.NewBuffer(nil)
				)
				err := copyWithOffset(r, w, tt.offset, tt.limit)
				require.NoError(t, err)
				require.Equal(t, tt.want, w.String())
			})
		}
	})

	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			name   string
			srs    string
			offset int64
			limit  int64
		}{
			// TODO: добавить тесты ошибок
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var (
					r = bytes.NewReader([]byte(tt.srs))
					w = bytes.NewBuffer(nil)
				)
				err := copyWithOffset(r, w, tt.offset, tt.limit)
				require.Error(t, err)
			})
		}
	})
}

func compareFiles(t *testing.T, f1, f2 string) {
	t.Helper()

	b1, err := os.ReadFile(f1)
	require.NoError(t, err)

	b2, err := os.ReadFile(f2)
	require.NoError(t, err)

	require.True(t, bytes.Equal(b1, b2))
}

func TestCopy(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		tests := []struct {
			offset int64
			limit  int64
		}{
			{0, 0},
			{0, 10},
			{0, 1000},
			{0, 10000},
			{100, 1000},
			{6000, 1000},
		}
		for _, tt := range tests {
			t.Run(fmt.Sprintf("offset=%d_limit=%d", tt.offset, tt.limit), func(t *testing.T) {
				err := Copy("testdata/input.txt", "/tmp/out.txt", tt.offset, tt.limit)
				require.NoError(t, err)

				compareFiles(t, fmt.Sprintf("testdata/out_offset%d_limit%d.txt", tt.offset, tt.limit), "/tmp/out.txt")
				os.Remove("/tmp/out.txt")
			})
		}
	})

	t.Run("errors: negative offset or limit", func(t *testing.T) {
		tests := []struct {
			name   string
			offset int64
			limit  int64
		}{
			{"negative offset", -1, 0},
			{"negative limit", 0, -1},
			{"offset exceeds file size", 10_000, 0},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := Copy("testdata/input.txt", "/tmp/out.txt", tt.offset, tt.limit)
				require.Error(t, err)
			})
		}
	})

	t.Run("errors: empty source or destination", func(t *testing.T) {
		t.Run("empty source name", func(t *testing.T) {
			err := Copy("", "/tmp/out.txt", 0, 0)
			require.Error(t, err)
		})

		t.Run("empty destination name", func(t *testing.T) {
			err := Copy("testdata/input.txt", "", 0, 0)
			require.Error(t, err)
		})

		t.Run("source file doesn't exist", func(t *testing.T) {
			err := Copy("testdata/input1.txt", "/tmp/out.txt", 0, 0)
			require.Error(t, err)
		})

		t.Run("no permission to read source", func(t *testing.T) {
			os.Chmod("testdata/no_permissios.txt", 0o020)

			err := Copy("testdata/no_permissios.txt", "/tmp/out.txt", 0, 0)
			require.Error(t, err)

			os.Chmod("testdata/no_permissios.txt", 0o420)
		})

		t.Run("no permission to write destination", func(t *testing.T) {
			err := Copy("testdata/input.txt", "/no_permissios.txt", 0, 0)
			require.Error(t, err)
		})

		t.Run("source is directory", func(t *testing.T) {
			err := Copy("testdata", "/tmp/out.txt", 0, 0)
			require.Error(t, err)
		})
	})

	t.Run("equal files", func(t *testing.T) {
		f1, err := os.Open("testdata/input.txt")
		require.NoError(t, err)
		defer f1.Close()

		s1, err := f1.Stat()
		require.NoError(t, err)

		f2, err := os.Open("testdata/input.txt")
		require.NoError(t, err)

		s2, err := f2.Stat()
		require.NoError(t, err)

		require.True(t, os.SameFile(s1, s2))
	})
}
