package hw09structvalidator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMin(t *testing.T) {
	t.Run("valid min", func(t *testing.T) {
		r, err := newMinValidator("10")
		require.NoError(t, err)
		require.NoError(t, r.Validate(10))
	})
	t.Run("invalid min", func(t *testing.T) {
		r, err := newMinValidator("10")
		require.NoError(t, err)
		require.Error(t, r.Validate(9))
	})
	t.Run("invalid tag", func(t *testing.T) {
		r, err := newMinValidator("a")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNotValidTag)
		require.Nil(t, r)
	})
}
