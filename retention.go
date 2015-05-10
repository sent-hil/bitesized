package bitesized

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type Retention map[string][]int

func (b *Bitesized) Retention(e string, from, till time.Time, i Interval) ([]Retention, error) {
	e = dasherize(e)
	retentions := []Retention{}

	start := till
	for {
		end := start
		counts := []int{}

		for {
			sKey := b.intervalkey(e, start, i)
			eKey := b.intervalkey(e, end, i)

			c, err := redis.Int(b.store.Do("BITOP", "AND", sKey+eKey, sKey, eKey))
			if err != nil {
				return nil, err
			}

			counts = append(counts, c)

			if _, err := b.store.Do("DEL", sKey+eKey); err != nil {
				return nil, err
			}

			if end = end.Add(-getDuration(i)); from.After(end) {
				break
			}
		}

		r := Retention{start.String(): counts}
		retentions = append(retentions, r)

		if start = start.Add(-getDuration(i)); from.After(start) {
			break
		}
	}

	return retentions, nil
}
