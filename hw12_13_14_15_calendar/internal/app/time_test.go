package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDay(t *testing.T) {
	t.Run("startOfDay", func(t *testing.T) {
		tests := []struct {
			name string
			t    time.Time
			want time.Time
		}{
			{
				"какое-то время",
				time.Date(2024, 5, 23, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 5, 23, 0, 0, 0, 0, time.UTC),
			},
			{
				"начало дня",
				time.Date(2024, 5, 23, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 5, 23, 0, 0, 0, 0, time.UTC),
			},
			{
				"конец дня",
				time.Date(2024, 5, 23, 23, 59, 59, 0, time.UTC),
				time.Date(2024, 5, 23, 0, 0, 0, 0, time.UTC),
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				require.Equal(t, tt.want, startOfDay(tt.t))
			})
		}
	})

	t.Run("endOfDay", func(t *testing.T) {
		tests := []struct {
			name string
			t    time.Time
			want time.Time
		}{
			{
				"обычный день",
				time.Date(2024, 5, 23, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 5, 23, 23, 59, 59, 0, time.UTC),
			},
			{
				"начало дня",
				time.Date(2024, 5, 23, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 5, 23, 23, 59, 59, 0, time.UTC),
			},
			{
				"конец дня",
				time.Date(2024, 5, 23, 23, 59, 59, 0, time.UTC),
				time.Date(2024, 5, 23, 23, 59, 59, 0, time.UTC),
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				require.Equal(t, tt.want, endOfDay(tt.t))
			})
		}
	})
}

func TestWeek(t *testing.T) {
	t.Run("startOfWeek", func(t *testing.T) {
		tests := []struct {
			name string
			t    time.Time
			want time.Time
		}{
			{
				"понедельник",
				time.Date(2024, 5, 20, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC),
			},
			{
				"вторник",
				time.Date(2024, 5, 21, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC),
			},
			{
				"среда",
				time.Date(2024, 5, 23, 22, 34, 56, 0, time.UTC),
				time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC),
			},
			{
				"четверг",
				time.Date(2024, 5, 23, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC),
			},
			{
				"пятница",
				time.Date(2024, 5, 24, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC),
			},
			{
				"суббота",
				time.Date(2024, 5, 25, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC),
			},
			{
				"воскресенье",
				time.Date(2024, 5, 26, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC),
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				require.Equal(t, tt.want, startOfWeek(tt.t))
			})
		}
	})

	t.Run("endOfWeek", func(t *testing.T) {
		tests := []struct {
			name string
			t    time.Time
			want time.Time
		}{
			{
				"понедельник",
				time.Date(2024, 5, 20, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 5, 26, 23, 59, 59, 0, time.UTC),
			},
			{
				"вторник",
				time.Date(2024, 5, 21, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 5, 26, 23, 59, 59, 0, time.UTC),
			},
			{
				"среда",
				time.Date(2024, 5, 22, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 5, 26, 23, 59, 59, 0, time.UTC),
			},
			{
				"четверг",
				time.Date(2024, 5, 23, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 5, 26, 23, 59, 59, 0, time.UTC),
			},
			{
				"пятница",
				time.Date(2024, 5, 24, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 5, 26, 23, 59, 59, 0, time.UTC),
			},
			{
				"суббота",
				time.Date(2024, 5, 25, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 5, 26, 23, 59, 59, 0, time.UTC),
			},
			{
				"воскресенье",
				time.Date(2024, 5, 26, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 5, 26, 23, 59, 59, 0, time.UTC),
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				require.Equal(t, tt.want, endOfWeek(tt.t))
			})
		}
	})
}

func TestMonth(t *testing.T) {
	t.Run("startOfMonth", func(t *testing.T) {
		tests := []struct {
			name string
			t    time.Time
			want time.Time
		}{
			{
				"обычный месяц",
				time.Date(2024, 5, 23, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				"начало месяца",
				time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				"конец месяца",
				time.Date(2024, 5, 31, 23, 59, 59, 0, time.UTC),
				time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				"февраль високосного года",
				time.Date(2024, 2, 29, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				require.Equal(t, tt.want, startOfMonth(tt.t))
			})
		}
	})

	t.Run("endOfMonth", func(t *testing.T) {
		tests := []struct {
			name string
			t    time.Time
			want time.Time
		}{
			{
				"обычный месяц",
				time.Date(2024, 5, 23, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 5, 31, 23, 59, 59, 0, time.UTC),
			},
			{
				"начало месяца",
				time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 5, 31, 23, 59, 59, 0, time.UTC),
			},
			{
				"конец месяца",
				time.Date(2024, 5, 31, 23, 59, 59, 0, time.UTC),
				time.Date(2024, 5, 31, 23, 59, 59, 0, time.UTC),
			},
			{
				"февраль високосного года",
				time.Date(2024, 2, 20, 12, 34, 56, 0, time.UTC),
				time.Date(2024, 2, 29, 23, 59, 59, 0, time.UTC),
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				require.Equal(t, tt.want, endOfMonth(tt.t))
			})
		}
	})
}
