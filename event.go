package bitesized

import (
	"regexp"
	"time"

	"github.com/garyburd/redigo/redis"
)

func (b *Bitesized) TrackEvent(evnt, user string, tstamp time.Time) error {
	return b.changeBit(evnt, user, tstamp, On)
}

func (b *Bitesized) UntrackEvent(evnt, user string, tstamp time.Time) error {
	return b.changeBit(evnt, user, tstamp, Off)
}

func (b *Bitesized) CountEvent(e string, t time.Time, i Interval) (int, error) {
	key := b.intervalkey(e, t, i)
	return redis.Int(b.store.Do("BITCOUNT", key))
}

func (b *Bitesized) DidEvent(e, u string, t time.Time, i Interval) (bool, error) {
	key := b.intervalkey(e, t, i)

	offset, err := b.getOrSetUser(u)
	if err != nil {
		return false, err
	}

	return redis.Bool(b.store.Do("GETBIT", key, offset))
}

func (b *Bitesized) GetEvents(prefix string) ([]string, error) {
	prefix = b.key(EventPrefixKey, prefix)
	allkeys, err := redis.Strings(b.store.Do("KEYS", prefix))
	if err != nil {
		return nil, err
	}

	rr := map[string]bool{}
	keys := []string{}

	for _, key := range allkeys {
		r := regexp.MustCompile(EventRegex)
		results := r.FindAllStringSubmatch(key, -1)

		if len(results) == 0 || len(results[0]) == 0 {
			continue
		}

		evnt := results[0][1]
		if _, ok := rr[evnt]; !ok {
			rr[evnt] = true
			keys = append(keys, key)
		}
	}

	return keys, nil
}
