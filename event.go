package bitesized

import (
	"regexp"
	"time"

	"github.com/garyburd/redigo/redis"
)

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

func (b *Bitesized) GetEvents(prefix string) ([]string, error) {
	prefix = b.key(EventPrefixKey, prefix)
	allkeys, err := redis.Strings(b.store.Do("KEYS", prefix))
	if err != nil {
		return nil, err
	}

	rr := map[string]bool{}
	keys := []string{}

	for _, k := range allkeys {
		r := regexp.MustCompile(EventRegex)
		results := r.FindAllStringSubmatch(k, -1)

		if len(results) == 0 || len(results[0]) == 0 {
			return keys, nil
		}

		evnt := results[0][1]
		if _, ok := rr[evnt]; !ok {
			rr[evnt] = true
			keys = append(keys, k)
		}
	}

	return keys, nil
}
