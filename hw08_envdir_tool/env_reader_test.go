package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadEnvFromBytes(t *testing.T) {
	t.Run("первая строка в файле - это значение переменной", func(t *testing.T) {
		tests := []struct {
			name  string
			input []byte
			want  string
		}{
			{"одна строка", []byte("test"), "test"},
			{"несколько строк", []byte("test1\ntest2"), "test1"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := readEnvFromBytes(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.want, got.Value)
			})
		}
	})

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

	t.Run("комплексный тест", func(t *testing.T) {
		test := []byte("a\000b \t\nc")
		got, err := readEnvFromBytes(test)
		require.NoError(t, err)
		require.Equal(t, "a\nb", got.Value)
	})
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
