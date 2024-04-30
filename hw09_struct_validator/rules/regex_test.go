package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegex(t *testing.T) {
	t.Run("valid regexp", func(t *testing.T) {
		r, err := newRegexValidator("^[a-z]+$")
		require.NoError(t, err)
		require.NoError(t, r.Validate("abcd"))
	})
	t.Run("invalid regexp", func(t *testing.T) {
		r, err := newRegexValidator("^[a-z]+$")
		require.NoError(t, err)
		require.Error(t, r.Validate("1234"))
	})
	t.Run("invalid tag", func(t *testing.T) {
		r, err := newRegexValidator("")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNotValidTag)
		require.Nil(t, r)
	})
}
