package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_copyWithOffset(t *testing.T) {
	t.Run("простые тесты", func(t *testing.T) {
		tests := []struct {
			name   string
			srs    string
			offset int64
			limit  int64
			want   string
		}{
			{"чтение всего файла", "1234567890", 0, 10, "1234567890"},
			{"чтение начала файла", "1234567890", 0, 5, "12345"},
			{"чтение середины файла", "1234567890", 3, 5, "45678"},
			{"чтение конца файла", "1234567890", 5, 5, "67890"},
			{"лимит больше размера файла без смещения", "1234567890", 0, 15, "1234567890"},
			{"лимит больше размера файла со смещением", "1234567890", 3, 15, "4567890"},
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
}

func TestCopy(t *testing.T) {
}
