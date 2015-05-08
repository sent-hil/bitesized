package bitesized

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
