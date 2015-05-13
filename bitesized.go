package bitesized

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// Bitesized is a client that can be used to track events and retrieve metrics.
type Bitesized struct {
	store redis.Conn

	Intervals []Interval

	// KeyPrefix is the prefix that'll be appended to all keys.
	KeyPrefix string
}

// NewClient initializes a Bitesized client with redis conn & default values.
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

func (b *Bitesized) getOrSetUser(user string) (int, error) {
	user = dasherize(user)

	script := redis.NewScript(3, getOrSetUserScript)
	raw, err := script.Do(b.store, b.userListKey(), user, b.userCounterKey())

	return redis.Int(raw, err)
}

func (b *Bitesized) storeIntervals(evnt string, offset int, t time.Time) error {
	b.store.Send("MULTI")

	for _, interval := range b.Intervals {
		key := b.intervalkey(evnt, t, interval)
		b.store.Send("SETBIT", key, offset, On)
	}

	_, err := b.store.Do("EXEC")

	return err
}
