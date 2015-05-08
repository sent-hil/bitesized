package bitesized

import (
	"github.com/garyburd/redigo/redis"
)

// Bitesized is a client that can be used to track events and retrieve metrics.
type Bitesized struct {
	store redis.Conn

	Intervals []Interval

	// KeyPrefix is the prefix that'll be appended to all keys.
	KeyPrefix string
}

// NewClient initializes a Bitesized client metrics. It initializes redis
// connection & default values for client.
func NewClient(redisuri string) (*Bitesized, error) {
	redissession, err := redis.Dial("tcp", redisuri)
	if err != nil {
		return nil, err
	}

	client := &Bitesized{
		store:     redissession,
		Intervals: DefaultIntervals,
		KeyPrefix: DefaultKeyPrefix,
	}

	return client, nil
}
