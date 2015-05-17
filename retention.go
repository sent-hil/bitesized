package bitesized

import "time"

type Retention map[string][]float64

func (b *Bitesized) Retention(e string, f, t time.Time, i Interval) ([]Retention, error) {
	if f.After(t) {
		return nil, ErrFromAfterTill
	}

	retentions := []Retention{}

	start := f
	for {
		end := start
		keyAggr := []string{}
		counts := []float64{}

		for {
			keyAggr = append(keyAggr, b.intervalkey(e, end, i))

			c, err := b.Operation(AND, keyAggr...)
			if err != nil {
				return nil, err
			}

			counts = append(counts, c)

			if end = end.Add(getDuration(end, i)); end.After(t) {
				break
			}
		}

		r := Retention{nearestInterval(start, i): counts}
		retentions = append(retentions, r)

		if start = start.Add(getDuration(start, i)); start.After(t) {
			break
		}
	}

	return retentions, nil
}

func (b *Bitesized) RetentionPercent(e string, f, t time.Time, i Interval) ([]Retention, error) {
	retentions, err := b.Retention(e, f, t, i)
	if err != nil {
		return nil, err
	}

	for _, rets := range retentions {
		for timekey, values := range rets {
			first := values[0]
			percents := []float64{first}

			for _, r := range values[1:] {
				var value float64
				if first != 0 {
					value = r / first
				}

				percents = append(percents, value*100)
			}

			rets[timekey] = percents
		}
	}

	return retentions, nil
}
