package rules_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trb1maker/otus_golang_home_work/hw09_struct_validator/rules"
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
				r, err := rules.NewStringRule(tt.tag)
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
			{"len", "len:5", "123456", rules.ErrLen},
			{"regex", "regexp:^\\d{5}$", "123456", rules.ErrRegexp},
			{"enum", "in:1,2,3", "4", rules.ErrEnum},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r, err := rules.NewStringRule(tt.tag)
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
				r, err := rules.NewStringRule(tt.tag)
				require.ErrorIs(t, err, rules.ErrNotValidTag)
				require.Nil(t, r)
			})
		}
	})
	t.Run("structed tag", func(t *testing.T) {
		tag := "len:1|regexp:^\\d$|in:1,2,3"
		r, err := rules.NewStringRule(tag)
		require.NoError(t, err)
		require.NoError(t, r.Validate("1"))
	})
	t.Run("structed tag with error", func(t *testing.T) {
		tag := "len:1|regexp:^\\d$|in:1,2,3"
		r, err := rules.NewStringRule(tag)
		require.NoError(t, err)
		require.ErrorIs(t, r.Validate("4"), rules.ErrEnum)
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
				r, err := rules.NewIntRule(tt.tag)
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
			{"min", "min:10", 9, rules.ErrMin},
			{"max", "max:10", 11, rules.ErrMax},
			{"enum", "in:1,2,3", 4, rules.ErrEnum},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r, err := rules.NewIntRule(tt.tag)
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
				r, err := rules.NewIntRule(tt.tag)
				require.ErrorIs(t, err, rules.ErrNotValidTag)
				require.Nil(t, r)
			})
		}
	})
	t.Run("structed tag", func(t *testing.T) {
		tag := "min:1|max:10|in:1,2,3"
		r, err := rules.NewIntRule(tag)
		require.NoError(t, err)
		require.NoError(t, r.Validate(1))
	})
	t.Run("structed tag with error", func(t *testing.T) {
		tag := "min:1max:10|in:1,2,3"
		r, err := rules.NewIntRule(tag)
		require.ErrorIs(t, err, rules.ErrNotValidTag)
		require.Nil(t, r)
	})
}
