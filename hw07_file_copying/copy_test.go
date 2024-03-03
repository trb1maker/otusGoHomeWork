package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_copyWithOffset(t *testing.T) {
	t.Run("простые тесты", func(t *testing.T) {
		type writeStringer interface {
			io.Writer
			fmt.Stringer
		}
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
					r io.ReaderAt   = bytes.NewReader([]byte(tt.srs))
					w writeStringer = bytes.NewBuffer(nil)
				)
				err := copyWithOffset(r, w, tt.offset, tt.limit)
				require.NoError(t, err)
				require.Equal(t, tt.want, w.String())
			})
		}
	})
}

func readExprected(t *testing.T, fileName string) fmt.Stringer {
	f, err := os.Open(fileName)
	if err != nil {
		t.Skip(err)
	}
	defer f.Close()
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, f)
	if err != nil {
		t.Skip(err)
	}
	return buf
}

func Test_copy(t *testing.T) {
	t.Run("стандартные тесты", func(t *testing.T) {
		tests := []struct {
			offset int64
			limit  int64
		}{
			{0, 0},
			{0, 10},
			{0, 1000},
			{0, 10000},
			{100, 1000},
			{6000, 10000},
		}
		for _, tt := range tests {
			fileName := fmt.Sprintf("./testdata/out_offset%d_limit%d.txt", tt.offset, tt.limit)
			t.Run(fileName, func(t *testing.T) {
				f, err := os.Open("./testdata/input.txt")
				if err != nil {
					t.Skip(err)
				}
				defer f.Close()
				buf := bytes.NewBuffer(nil)
				err = copy(f, buf, tt.offset, tt.limit)
				require.NoError(t, err)
				require.Equal(t, readExprected(t, fileName).String(), buf.String())
			})
		}
	})
}

func TestCopy(t *testing.T) {
	t.Run("чтение и запись без ошибок", func(t *testing.T) {

	})

	t.Run("чтение и запись с ошибками", func(t *testing.T) {

	})
}
