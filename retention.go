package bitesized

import "time"

type Retention map[string][]int

func (b *Bitesized) Retention(e string, f, t time.Time, i Interval) ([]Retention, error) {
	if f.After(t) {
		return nil, ErrFromAfterTill
	}

	retentions := []Retention{}

	start := f
	for {
		end := start
		keyAggr := []string{}
		counts := []int{}

		for {
			keyAggr = append(keyAggr, b.intervalkey(e, end, i))

			c, err := b.Operation(AND, keyAggr...)
			if err != nil {
				return nil, err
			}

			counts = append(counts, c)

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
