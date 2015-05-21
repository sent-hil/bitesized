package bitesized

import (
	"time"

	"github.com/jinzhu/now"
)

// Interval define which time intervals to track events. Ex: `Month` interval
// turns on bit for that user in the specified month's bit array. Multiple
// intervals can be selected.
type Interval int

const (
	All Interval = iota
	Hour
	Day
	Week
	Month
	Year
)

func nearestInterval(t time.Time, interval Interval) string {
	n := now.New(t.UTC())

	switch interval {
	case All:
		return "all"
	case Day:
		layout := "day:2006-01-02"
		return n.BeginningOfDay().Format(layout)
	case Week:
		layout := "week:2006-01-02"
		return n.BeginningOfWeek().Format(layout)
	case Month:
		layout := "month:2006-01"
		return n.BeginningOfMonth().Format(layout)
	case Year:
		layout := "year:2006"
		return n.BeginningOfYear().Format(layout)
	}

	layout := "hour:2006-01-02-15:04"
	return n.BeginningOfHour().Format(layout)
}

func getDuration(t time.Time, i Interval) time.Duration {
	switch i {
	case Day:
		return 24 * time.Hour
	case Week:
		return 7 * 24 * time.Hour
	case Month:
		noOfDays := daysIn(t.Month(), t.Year())
		return time.Duration(noOfDays) * 24 * time.Hour
	case Year:
		return 365 * 24 * time.Hour
	}

	return time.Hour
}

func daysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
