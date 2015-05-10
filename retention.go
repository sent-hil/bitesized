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

	start := t
	for {
		end := start
		counts := []int{}

		for {
			sKey := b.intervalkey(e, start, i)
			eKey := b.intervalkey(e, end, i)
			rKey := randSeq(5)

			if _, err := b.store.Do("BITOP", "AND", rKey, sKey, eKey); err != nil {
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

			if end = end.Add(-getDuration(i)); f.After(end) {
				break
			}
		}

		r := Retention{start.String(): counts}
		retentions = append(retentions, r)

		if start = start.Add(-getDuration(i)); f.After(start) {
			break
		}
	}

	return retentions, nil
}
