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
	Hour Interval = iota
	Day
	Week
	Month
	Quater
	Year
	Decade
)

func nearestInterval(t time.Time, interval Interval) string {
	n := now.New(t.UTC())

	switch interval {
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

	layout := "hour:2006-01-02 15:04"
	return n.BeginningOfHour().Format(layout)
}
