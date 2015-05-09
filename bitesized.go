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

func (b *Bitesized) TrackEvent(name, username string, tstamp time.Time) error {
	if name == "" || username == "" {
		return ErrInvalidArg
	}

	offset, err := b.getOrSetUser(username)
	if err != nil {
		return err
	}

	return b.storeIntervals(name, offset, tstamp)
}

func (b *Bitesized) getOrSetUser(username string) (int, error) {
	script := redis.NewScript(3, getOrsetUserScript)
	raw, err := script.Do(b.store, b.userListKey(), username, b.userCounterKey())

	return redis.Int(raw, err)
}

func (b *Bitesized) storeIntervals(name string, offset int, tstamp time.Time) error {
	b.store.Send("MULTI")

	for _, interval := range b.Intervals {
		key := b.intervalKey(tstamp, interval)
		b.store.Send("SETBIT", key, offset, On)
	}

	_, err := b.store.Do("EXEC")

	return err
}

func (b *Bitesized) intervalKey(tstamp time.Time, interval Interval) string {
	intervalstr := nearestInterval(tstamp, interval)
	return b.key(intervalstr)
}

func (b *Bitesized) userListKey() string {
	return b.key(UserListKey)
}

func (b *Bitesized) userCounterKey() string {
	return b.key(UserCounterKey)
}

func (b *Bitesized) key(suffix string) string {
	key := suffix

	if b.KeyPrefix != "" {
		key = b.KeyPrefix + ":" + key
	}

	return key
}
