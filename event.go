package bitesized

import "time"

type Event struct {
	Name      string
	Username  string
	Timestamp time.Time
}
