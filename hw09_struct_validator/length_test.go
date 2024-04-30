package hw09structvalidator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLength(t *testing.T) {
	t.Run("valid length", func(t *testing.T) {
		r, err := newLengthValidator("4")
		require.NoError(t, err)
		require.NoError(t, r.Validate("abcd"))
	})
	t.Run("invalid length", func(t *testing.T) {
		r, err := newLengthValidator("10")
		require.NoError(t, err)
		require.Error(t, r.Validate("abcd"))
	})
	t.Run("invalid tag", func(t *testing.T) {
		r, err := newLengthValidator("a")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNotValidTag)
		require.Nil(t, r)
	})
}
