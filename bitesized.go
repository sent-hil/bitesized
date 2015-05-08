package bitesized

import (
	"github.com/garyburd/redigo/redis"
)

type Bitesized struct {
	store     redis.Conn
	Intervals []Interval
	KeyPrefix string
}

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
