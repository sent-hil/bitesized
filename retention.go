package bitesized

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type Retention map[string][]int

func (b *Bitesized) Retention(e string, f, t time.Time, i Interval) ([]Retention, error) {
	if f.After(t) {
		return nil, ErrFromAfterTill
	}

	e = dasherize(e)
	retentions := []Retention{}

	start := f
	for {
		end := start
		eKeys := []interface{}{}
		counts := []int{}

		for {
			rKey := randSeq(20)

			eKey := b.intervalkey(e, end, i)
			eKeys = append(eKeys, eKey)

			args := []interface{}{"AND", rKey}
			args = append(args, eKeys...)

			// TODO: use lua scripting
			if _, err := b.store.Do("BITOP", args...); err != nil {
				return nil, err
			}

			c, err := redis.Int(b.store.Do("BITCOUNT", rKey))
			if err != nil {
				return nil, err
			}

			counts = append(counts, c)

			if _, err := b.store.Do("DEL", rKey); err != nil {
				return nil, err
			}

			if end = end.Add(getDuration(i)); end.After(t) {
				break
			}
		}

		r := Retention{nearestInterval(start, i): counts}
		retentions = append(retentions, r)

		if start = start.Add(getDuration(i)); start.After(t) {
			break
		}
	}

	return retentions, nil
}
