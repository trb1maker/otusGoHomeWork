package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMin(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		require.NoError(t, validateMin(int64(10), "10"))
		require.NoError(t, validateMin(uint64(11), "10"))
		require.NoError(t, validateMin(float64(10), "7.5"))
	})
	t.Run("base error", func(t *testing.T) {
		require.ErrorIs(t, validateMin(int64(10), "11"), ErrMin)
		require.ErrorIs(t, validateMin(uint64(10), "11"), ErrMin)
		require.ErrorIs(t, validateMin(float64(10), "11"), ErrMin)
	})
}

func TestMax(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		require.NoError(t, validateMax(int64(-10), "10"))
		require.NoError(t, validateMax(uint64(11), "12"))
		require.NoError(t, validateMax(float64(10), "11.3"))
	})
	t.Run("base error", func(t *testing.T) {
		require.ErrorIs(t, validateMax(int64(12), "-11"), ErrMax)
		require.ErrorIs(t, validateMax(uint64(12), "11"), ErrMax)
		require.ErrorIs(t, validateMax(float64(12), "11.9"), ErrMax)
	})
}

func TestLen(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		require.NoError(t, validateLen("10", "2"))
	})
	t.Run("base error", func(t *testing.T) {
		require.ErrorIs(t, validateLen("10", "11"), ErrLength)
	})
}

func TestIn(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		require.NoError(t, validateIn("a", "b,a,c"))
		require.NoError(t, validateIn(int64(-1), "-1,2,3"))
		require.NoError(t, validateIn(uint64(2), "2,3,4"))
		require.NoError(t, validateIn(float64(4.5), "2.5,3.5,4.5"))
	})
	t.Run("base error", func(t *testing.T) {
		require.ErrorIs(t, validateIn("a", "b,c,d"), ErrIn)
	})
}

func TestRegexp(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		require.NoError(t, validateRegexp("test", "^[a-z]{4}$"))
	})
	t.Run("base error", func(t *testing.T) {
		require.ErrorIs(t, validateRegexp("test", "^[a-z]{5}$"), ErrRegexp)
	})
}

func TestNewRule(t *testing.T) {
	r, err := newRule("len:10")
	require.NoError(t, err)
	require.Equal(t, "len", r.name)
	require.Equal(t, "10", r.value)

	_, err = newRule("len")
	require.ErrorIs(t, err, ErrInvalidRule)

	_, err = newRule("")
	require.ErrorIs(t, err, ErrInvalidRule)
}

func TestParseRules(t *testing.T) {
	rules, err := parseRules("len:10|in:a|regexp:^[a-z]{4}$")
	require.NoError(t, err)
	require.Equal(t, 3, len(rules))
	require.Equal(t, "len", rules[0].name)

	_, err = parseRules("")
	require.ErrorIs(t, err, ErrInvalidRule)

	_, err = parseRules("len:10||in:a|regexp:^[a-z]{4}$")
	require.ErrorIs(t, err, ErrInvalidRule)
}

func TestValidateString(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		tag := "len:1|in:a,b,c"
		rr, err := parseRules(tag)
		require.NoError(t, err)

		require.NoError(t, validateString("c", rr))
		require.NoError(t, ValidateString([]string{"b"}, tag))
	})
}

func TestValidateNumber(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		tag := "min:10|max:20"
		rr, err := parseRules(tag)
		require.NoError(t, err)

		require.NoError(t, validateNumber(int64(12), rr))
		require.NoError(t, ValidateNumber([]int64{11}, tag))
	})
}
