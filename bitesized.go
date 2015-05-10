package bitesized

import (
	"strings"
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

func (b *Bitesized) TrackEvent(name, user string, tstamp time.Time) error {
	if name == "" || user == "" {
		return ErrInvalidArg
	}

	name = dasherize(name)
	user = dasherize(user)

	offset, err := b.getOrSetUser(user)
	if err != nil {
		return err
	}

	return b.storeIntervals(name, offset, tstamp)
}

func (b *Bitesized) CountEvent(n string, t time.Time, i Interval) (int, error) {
	n = dasherize(n)
	key := b.intervalkey(n, t, i)

	return redis.Int(b.store.Do("BITCOUNT", key))
}

func (b *Bitesized) DidEvent(n, u string, t time.Time, i Interval) (bool, error) {
	n = dasherize(n)
	u = dasherize(u)

	key := b.intervalkey(n, t, i)

	offset, err := b.getOrSetUser(u)
	if err != nil {
		return false, err
	}

	return redis.Bool(b.store.Do("GETBIT", key, offset))
}

func (b *Bitesized) getOrSetUser(user string) (int, error) {
	script := redis.NewScript(3, getOrsetUserScript)
	raw, err := script.Do(b.store, b.userListKey(), user, b.userCounterKey())

	return redis.Int(raw, err)
}

func (b *Bitesized) storeIntervals(name string, offset int, t time.Time) error {
	b.store.Send("MULTI")

	for _, interval := range b.Intervals {
		key := b.intervalkey(name, t, interval)
		b.store.Send("SETBIT", key, offset, On)
	}

	_, err := b.store.Do("EXEC")

	return err
}

func (b *Bitesized) intervalkey(name string, t time.Time, i Interval) string {
	intervalkey := nearestInterval(t, i)
	return b.key(name, intervalkey)
}

func (b *Bitesized) userListKey() string {
	return b.key(UserListKey)
}

func (b *Bitesized) userCounterKey() string {
	return b.key(UserCounterKey)
}

func (b *Bitesized) key(suffix ...string) string {
	key := strings.Join(suffix, ":")

	if b.KeyPrefix != "" {
		key = b.KeyPrefix + ":" + key
	}

	return key
}

func dasherize(name string) string {
	return strings.Join(strings.Split(name, " "), "-")
}
