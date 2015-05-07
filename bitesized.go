package bitesized

import (
	"github.com/garyburd/redigo/redis"
)

type Bitesized struct {
	Store     redis.Conn
	Intervals []Interval
	KeyPrefix string
}

func NewClient(r redis.Conn) *Bitesized {
	return &Bitesized{
		Store:     r,
		Intervals: DefaultIntervals,
		KeyPrefix: DefaultKeyPrefix,
	}
}
