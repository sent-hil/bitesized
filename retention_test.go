package bitesized

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRetention(t *testing.T) {
	Convey("", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		client.Intervals = []Interval{Hour}

		times := []time.Time{
			time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2015, time.January, 1, 1, 0, 0, 0, time.UTC),
			time.Date(2015, time.January, 1, 2, 0, 0, 0, time.UTC),
		}

		from, till := times[0], times[len(times)-1]

		Convey("It should return result with empty values", func() {
			retention, err := client.Retention("dodge rock", from, till, Hour)
			So(err, ShouldBeNil)

			So(len(retention), ShouldEqual, 3)

			for _, counts := range retention[0] {
				So(counts, ShouldContain, 0)
				So(counts, ShouldNotContain, 1)
			}
		})

		Convey("It should return result with values", func() {
			for _, t := range times {
				err := client.TrackEvent("dodge rock", user, t)
				So(err, ShouldBeNil)
			}

			retention, err := client.Retention("dodge rock", from, till, Hour)
			So(err, ShouldBeNil)

			So(len(retention), ShouldEqual, 3)

			for _, counts := range retention[0] {
				So(counts, ShouldContain, 1)
				So(counts, ShouldNotContain, 0)
			}

			Reset(func() { client.store.Do("FLUSHALL") })
		})
	})
}
