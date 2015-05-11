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

func (b *Bitesized) TrackEvent(evnt, user string, tstamp time.Time) error {
	if evnt == "" || user == "" {
		return ErrInvalidArg
	}

	evnt = dasherize(evnt)
	user = dasherize(user)

	offset, err := b.getOrSetUser(user)
	if err != nil {
		return err
	}

	return b.storeIntervals(evnt, offset, tstamp)
}

func (b *Bitesized) CountEvent(e string, t time.Time, i Interval) (int, error) {
	e = dasherize(e)
	key := b.intervalkey(e, t, i)

	return redis.Int(b.store.Do("BITCOUNT", key))
}

func (b *Bitesized) DidEvent(e, u string, t time.Time, i Interval) (bool, error) {
	e = dasherize(e)
	u = dasherize(u)

	key := b.intervalkey(e, t, i)

	offset, err := b.getOrSetUser(u)
	if err != nil {
		return false, err
	}

	return redis.Bool(b.store.Do("GETBIT", key, offset))
}

func (b *Bitesized) getOrSetUser(user string) (int, error) {
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
