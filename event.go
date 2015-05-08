package bitesized

import "time"

// Event is an action taken by a single user at a particular time.
type Event struct {
	Name      string
	Username  string
	Timestamp time.Time
}
