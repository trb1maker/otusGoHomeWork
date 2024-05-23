package app

import "time"

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func endOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}

func startOfWeek(t time.Time) time.Time {
	weekday := t.Weekday()
	if weekday == 0 {
		weekday = 7
	}
	return time.Date(t.Year(), t.Month(), t.Day()-int(weekday)+1, 0, 0, 0, 0, t.Location())
}

func endOfWeek(t time.Time) time.Time {
	weekday := t.Weekday()
	if weekday == 0 {
		weekday = 7
	}
	return time.Date(t.Year(), t.Month(), t.Day()-int(weekday)+7, 23, 59, 59, 0, t.Location())
}

func startOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func endOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 1, 23, 59, 59, 0, t.Location()).AddDate(0, 0, -1)
}
