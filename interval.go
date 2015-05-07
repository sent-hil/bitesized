package bitesized

type Interval int

const (
	Hour Interval = iota
	Day
	Week
	Month
	Quater
	Year
)
