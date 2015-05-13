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

func (b *Bitesized) Operation(op Op, keys ...string) (int, error) {
	if op == NOT && len(keys) != 1 {
		return 0, ErrNotOpAcceptsOnekey
	}

	rKey := randSeq(20)

	args := []interface{}{op, rKey}
	for _, key := range keys {
		args = append(args, key)
	}

	if _, err := b.store.Do("BITOP", args...); err != nil {
		return 0, err
	}

	count, err := redis.Int(b.store.Do("BITCOUNT", rKey))
	if err != nil {
		return 0, err
	}

	if _, err := b.store.Do("DEL", rKey); err != nil {
		return 0, err
	}

	return count, nil
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
