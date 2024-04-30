package hw09structvalidator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnumInt(t *testing.T) {
	t.Run("valid enum", func(t *testing.T) {
		r, err := newEnumIntValidator("1,2,3")
		require.NoError(t, err)
		require.NoError(t, r.Validate(1))
	})
	t.Run("invalid enum", func(t *testing.T) {
		r, err := newEnumIntValidator("1,2,3")
		require.NoError(t, err)
		require.Error(t, r.Validate(4))
	})
	t.Run("invalid tag", func(t *testing.T) {
		r, err := newEnumIntValidator("a,b,c")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNotValidTag)
		require.Nil(t, r)
	})
}

func TestEnumString(t *testing.T) {
	t.Run("valid enum", func(t *testing.T) {
		r, err := newEnumStringValidator("a,b,c")
		require.NoError(t, err)
		require.NoError(t, r.Validate("a"))
	})
	t.Run("invalid enum", func(t *testing.T) {
		r, err := newEnumStringValidator("a,b,c")
		require.NoError(t, err)
		require.Error(t, r.Validate("d"))
	})
	t.Run("invalid tag", func(t *testing.T) {
		r, err := newEnumStringValidator("")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNotValidTag)
		require.Nil(t, r)
	})
}
