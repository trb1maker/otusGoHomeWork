package hw09structvalidator

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trb1maker/otus_golang_home_work/hw09_struct_validator/rules"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Age struct {
		Age int `validate:"min:18|max:50"`
		p   int `validate:"in:0,1"`
	}

	OtherAge struct {
		Age int `validate:"min:a|max:50"`
	}

	Token struct {
		Header    []byte `validate:"len:10"`
		Payload   []byte
		Signature []byte
	}

	Email struct {
		Email string `validate:"min:10"`
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	t.Run("нормальное валидирование", func(t *testing.T) {
		tests := []struct {
			name string
			in   interface{}
		}{
			{
				"user",
				User{
					ID:     "332207b5-2b11-4a0f-988d-adc73a025b86",
					Age:    20,
					Email:  "qGwGQ@example.com",
					Role:   "admin",
					Phones: []string{"12345678901", "12345678902"},
				},
			},
			{"response", Response{Code: 200, Body: "body"}},
			{"пустая структура", struct{}{}},
		}
		for _, tt := range tests {
			// tt := tt
			// t.Parallel()
			t.Run(tt.name, func(t *testing.T) {
				require.NoError(t, Validate(tt.in))
			})
		}
	})
	t.Run("ошибка валидации", func(t *testing.T) {
		tests := []struct {
			name     string
			in       interface{}
			expected error
		}{
			{"len", App{Version: "10.0"}, rules.ErrLen},
			{"in", Response{Code: 100, Body: "body"}, rules.ErrIn},
			{"нарушение одного из правил", Age{Age: 100}, rules.ErrMax},
		}
		for _, tt := range tests {
			// tt := tt
			// t.Parallel()
			t.Run(tt.name, func(t *testing.T) {
				err := Validate(tt.in)
				require.Error(t, err)
				require.ErrorAs(t, err, &tt.expected)
			})
		}
	})
	t.Run("ошибки в тэгах", func(t *testing.T) {
		tests := []struct {
			name     string
			in       interface{}
			expected error
		}{
			{"age", OtherAge{Age: 37}, rules.ErrInvalidRule},
			{"email", Email{Email: "qGwGQ@example.com"}, rules.ErrInvalidRule},
			// Тест работает, но на него линтер ругается
			// File is not `gofumpt`-ed (gofumpt)
			// {
			// 	"token",
			// 	Token{
			// 		Header:    []byte("1234567890"),
			// 		Payload:   []byte("1234567890"),
			// 		Signature: []byte("1234567890"),
			// 	},
			// 	rules.ErrUnsupportedType,
			// },
			{"not a struct", 1, ErrNotStruct},
		}
		for _, tt := range tests {
			// tt := tt
			// t.Parallel()
			t.Run(tt.name, func(t *testing.T) {
				err := Validate(tt.in)
				require.Error(t, err)
				require.ErrorAs(t, err, &tt.expected)
			})
		}
	})
	t.Run("игнорирование приватных полей", func(t *testing.T) {
		// Значение p соответствует тэгу
		a := Age{Age: 34, p: 1}
		require.NoError(t, Validate(a))

		// Значение p тэгу не соответствует
		a = Age{Age: 34, p: 2}
		require.NoError(t, Validate(a))
	})
}
