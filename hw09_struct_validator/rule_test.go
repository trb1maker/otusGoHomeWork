package hw09structvalidator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	t.Run("valid string", func(t *testing.T) {
		tests := []struct {
			name  string
			tag   string
			value string
		}{
			{"len", "len:5", "12345"},
			{"regex", "regexp:^\\d{5}$", "12345"},
			{"enum", "in:1,2,3", "1"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r, err := NewStringRule(tt.tag)
				require.NoError(t, err)
				require.NoError(t, r.Validate(tt.value))
			})
		}
	})
	t.Run("invalid string", func(t *testing.T) {
		tests := []struct {
			name  string
			tag   string
			value string
			err   error
		}{
			{"len", "len:5", "123456", ErrLen},
			{"regex", "regexp:^\\d{5}$", "123456", ErrRegexp},
			{"enum", "in:1,2,3", "4", ErrEnum},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r, err := NewStringRule(tt.tag)
				require.NoError(t, err)
				require.ErrorIs(t, r.Validate(tt.value), tt.err)
			})
		}
	})
	t.Run("invalid tag", func(t *testing.T) {
		tests := []struct {
			name  string
			tag   string
			value string
		}{
			{"len", "len", "12345"},
			{"regexp", "", "12345"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r, err := NewStringRule(tt.tag)
				require.ErrorIs(t, err, ErrNotValidTag)
				require.Nil(t, r)
			})
		}
	})
	t.Run("structed tag", func(t *testing.T) {
		tag := "len:1|regexp:^\\d$|in:1,2,3"
		r, err := NewStringRule(tag)
		require.NoError(t, err)
		require.NoError(t, r.Validate("1"))
	})
	t.Run("structed tag with error", func(t *testing.T) {
		tag := "len:1|regexp:^\\d$|in:1,2,3"
		r, err := NewStringRule(tag)
		require.NoError(t, err)
		require.ErrorIs(t, r.Validate("4"), ErrEnum)
	})
}

func TestInt(t *testing.T) {
	t.Run("valid int", func(t *testing.T) {
		tests := []struct {
			name  string
			tag   string
			value int
		}{
			{"min", "min:10", 10},
			{"max", "max:10", 10},
			{"enum", "in:1,2,3", 1},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r, err := NewIntRule(tt.tag)
				require.NoError(t, err)
				require.NoError(t, r.Validate(tt.value))
			})
		}
	})
	t.Run("invalid int", func(t *testing.T) {
		tests := []struct {
			name  string
			tag   string
			value int
			err   error
		}{
			{"min", "min:10", 9, ErrMin},
			{"max", "max:10", 11, ErrMax},
			{"enum", "in:1,2,3", 4, ErrEnum},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r, err := NewIntRule(tt.tag)
				require.NoError(t, err)
				require.ErrorIs(t, r.Validate(tt.value), tt.err)
			})
		}
	})
	t.Run("invalid tag", func(t *testing.T) {
		tests := []struct {
			name  string
			tag   string
			value int
		}{
			{"min", "min", 10},
			{"max", "", 10},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r, err := NewIntRule(tt.tag)
				require.ErrorIs(t, err, ErrNotValidTag)
				require.Nil(t, r)
			})
		}
	})
	t.Run("structed tag", func(t *testing.T) {
		tag := "min:1|max:10|in:1,2,3"
		r, err := NewIntRule(tag)
		require.NoError(t, err)
		require.NoError(t, r.Validate(1))
	})
	t.Run("structed tag with error", func(t *testing.T) {
		tag := "min:1max:10|in:1,2,3"
		r, err := NewIntRule(tag)
		require.ErrorIs(t, err, ErrNotValidTag)
		require.Nil(t, r)
	})
}
