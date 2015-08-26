package bitesized

import (
	"math"
	"time"

	"github.com/jinzhu/now"
)

// Interval define which time intervals to track events. Ex: `Month` interval
// turns on bit for that user in the specified month's bit array. Multiple
// intervals can be selected.
type Interval int

const (
	All Interval = iota
	TenMinutes
	ThirtyMinutes
	Hour
	Day
	Biweekly
	Week
	Bimonthly
	Month
	Quarter
	Year
)

func handleMinuteInterval(t time.Time, n *now.Now, cycleLength int, keyName string) string {
	layout := keyName + ":2006-01-02-15:04"
	offset := t.Sub(n.BeginningOfHour())
	cycle := int(math.Floor(offset.Minutes() / float64(cycleLength)))
	return n.BeginningOfHour().Add(time.Duration(cycle*cycleLength) * time.Minute).Format(layout)
}

func nearestInterval(t time.Time, interval Interval) string {
	n := now.New(t.UTC())

	switch interval {
	case All:
		return "all"
	case TenMinutes:
		return handleMinuteInterval(t, n, 10, "ten_minutes")
	case ThirtyMinutes:
		return handleMinuteInterval(t, n, 30, "thirty_minutes")
	case Day:
		layout := "day:2006-01-02"
		return n.BeginningOfDay().Format(layout)
	case Biweekly:
		layout := "biweekly:2006-01-02"
		date := n.BeginningOfWeek()
		if offset := t.Sub(n.BeginningOfWeek()); offset.Hours() > 84 {
			date = date.Add(84 * time.Hour)
		}
		return date.Format(layout)
	case Week:
		layout := "week:2006-01-02"
		return n.BeginningOfWeek().Format(layout)
	case Bimonthly:
		layout := "bimonthly:2006-01-02"
		monthMiddle := n.EndOfMonth().Sub(n.BeginningOfMonth()) / 2
		date := n.BeginningOfMonth()
		if offset := t.Sub(n.BeginningOfMonth()); offset > monthMiddle {
			date = date.Add(monthMiddle)
		}
		return date.Format(layout)
	case Month:
		layout := "month:2006-01"
		return n.BeginningOfMonth().Format(layout)
	case Quarter:
		layout := "quarter:2006-01"
		return n.BeginningOfQuarter().Format(layout)
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
