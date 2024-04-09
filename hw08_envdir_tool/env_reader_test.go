package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadEnvFromBytes(t *testing.T) {
	t.Run("пробелы и символы табуляции в конце значения переменной должны быть удалены", func(t *testing.T) {
		tests := []struct {
			name  string
			input []byte
			want  string
		}{
			{"пробелы", []byte("test "), "test"},
			{"табуляции", []byte("test\t"), "test"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := readEnvFromBytes(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.want, got.Value)
			})
		}
	})

	t.Run("терминальные нули должны быть трансформированы в символ перевода строки", func(t *testing.T) {
		tests := []struct {
			name  string
			input []byte
			want  string
		}{
			{"терминальные нули", []byte("test\000"), "test\n"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := readEnvFromBytes(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.want, got.Value)
			})
		}
	})
}

func TestReadEnvFromFile(t *testing.T) {
	f := bytes.NewBufferString("test\ntest1")

	buf := make([]byte, 30, 120)
	offset := 0

	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			buf = buf[:offset]
			break
		}
		require.NoError(t, err)

		if i := bytes.IndexByte(buf[:offset+n], '\n'); i >= 0 {
			buf = buf[:offset+i]
			break
		}
		offset += n
	}
	require.Equal(t, "test", string(buf))
}

func TestReadDir(t *testing.T) {
	env, err := ReadDir("testdata/env")
	require.NoError(t, err)
	test := Environment{
		"BAR":   EnvValue{Value: "bar", NeedRemove: false},
		"EMPTY": EnvValue{Value: "", NeedRemove: false},
		"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
		"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
		"UNSET": EnvValue{Value: "", NeedRemove: true},
	}
	require.Equal(t, test, env)
}
