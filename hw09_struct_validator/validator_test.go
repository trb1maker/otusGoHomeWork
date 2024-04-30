package hw09structvalidator

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^[\\w\\d._-]+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	t.Run("valid users", func(t *testing.T) {
		f, err := os.Open("testdata/users.json")
		if err != nil {
			t.Skip(err)
		}
		defer f.Close()
		users := make([]User, 0, 25)
		if err := json.NewDecoder(f).Decode(&users); err != nil {
			t.Skip(err)
		}
		for _, u := range users {
			t.Run(u.ID, func(t *testing.T) {
				require.NoError(t, Validate(u))
			})
		}
	})
	t.Run("invalid users", func(t *testing.T) {
		u := User{
			ID:     "123",
			Name:   "John",
			Age:    32,
			Email:  "john@hw.ru",
			Role:   "admin",
			Phones: []string{"12345678901"},
		}
		require.Error(t, Validate(u))
	})
}
