package hw09structvalidator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMax(t *testing.T) {
	t.Run("valid max", func(t *testing.T) {
		r, err := newMaxValidator("10")
		require.NoError(t, err)
		require.NoError(t, r.Validate(10))
	})
	t.Run("invalid max", func(t *testing.T) {
		r, err := newMaxValidator("10")
		require.NoError(t, err)
		require.Error(t, r.Validate(11))
	})
	t.Run("invalid tag", func(t *testing.T) {
		r, err := newMaxValidator("10.1")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNotValidTag)
		require.Nil(t, r)
	})
}
