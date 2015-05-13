package bitesized

import (
	"strings"
	"time"
)

var (
	On  = 1
	Off = 0
)

var (
	EventPrefixKey = "event"
	UserListKey    = "user-list"
	UserCounterKey = "user-counter"
)

func (b *Bitesized) intervalkey(evnt string, t time.Time, i Interval) string {
	intervalkey := nearestInterval(t, i)
	return b.key(EventPrefixKey, evnt, intervalkey)
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
