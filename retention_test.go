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

		from := time.Date(2015, time.January, 2, 0, 1, 0, 0, time.UTC)
		till := from.Add(5 * time.Hour)

		data := map[string][]time.Time{
			"user1": []time.Time{
				from,
				from.Add(1 * time.Hour),
				from.Add(2 * time.Hour),
				from.Add(4 * time.Hour),
				from.Add(5 * time.Hour),
			},
			"user2": []time.Time{
				from.Add(1 * time.Hour),
			},
			"user3": []time.Time{
				from.Add(4 * time.Hour),
				from.Add(5 * time.Hour),
			},
		}

		Convey("It should error if from is after till", func() {
			till := from.Add(1 * time.Hour)

			_, err := client.Retention("dodge rock", till, from, Hour, 1)
			So(err, ShouldEqual, ErrFromAfterTill)
		})

		Convey("It should return result with empty values", func() {
			retention, err := client.Retention("dodge rock", from, till, Hour, 5)
			So(err, ShouldBeNil)

			So(len(retention), ShouldEqual, 6)

			for _, counts := range retention[0] {
				So(counts, ShouldContain, 0)
				So(counts, ShouldNotContain, 1)
			}
		})

		Convey("It should return result with values", func() {
			for user, times := range data {
				for _, t := range times {
					err := client.TrackEvent("dodge rock", user, t)
					So(err, ShouldBeNil)
				}
			}

			retention, err := client.Retention("dodge rock", from, till, Hour, 5)
			So(err, ShouldBeNil)

			So(len(retention), ShouldEqual, 6)

			for _, counts := range retention[0] {
				So(len(counts), ShouldEqual, 5)

				So(counts[0], ShouldEqual, 1)
				So(counts[1], ShouldEqual, 1)
				So(counts[2], ShouldEqual, 1)
				So(counts[3], ShouldEqual, 0)
				So(counts[4], ShouldEqual, 0)
			}

			for _, counts := range retention[1] {
				So(len(counts), ShouldEqual, 5)

				So(counts[0], ShouldEqual, 2)
				So(counts[1], ShouldEqual, 1)
				So(counts[2], ShouldEqual, 0)
				So(counts[3], ShouldEqual, 0)
				So(counts[4], ShouldEqual, 0)
			}

			for _, counts := range retention[1] {
				So(len(counts), ShouldEqual, 5)

				So(counts[0], ShouldEqual, 2)
				So(counts[1], ShouldEqual, 1)
				So(counts[2], ShouldEqual, 0)
				So(counts[3], ShouldEqual, 0)
				So(counts[4], ShouldEqual, 0)
			}

			for _, counts := range retention[2] {
				So(len(counts), ShouldEqual, 4)

				So(counts[0], ShouldEqual, 1)
				So(counts[1], ShouldEqual, 0)
				So(counts[2], ShouldEqual, 0)
				So(counts[3], ShouldEqual, 0)
			}

			for _, counts := range retention[3] {
				So(len(counts), ShouldEqual, 3)

				So(counts[0], ShouldEqual, 0)
				So(counts[1], ShouldEqual, 0)
				So(counts[2], ShouldEqual, 0)
			}

			for _, counts := range retention[4] {
				So(len(counts), ShouldEqual, 2)

				So(counts[0], ShouldEqual, 2)
				So(counts[1], ShouldEqual, 2)
			}

			for _, counts := range retention[5] {
				So(len(counts), ShouldEqual, 1)

				So(counts[0], ShouldEqual, 2)
			}

			Convey("It should results as percentages", func() {
				retention, err := client.RetentionPercent("dodge rock", from, till, Hour, 5)
				So(err, ShouldBeNil)

				for _, counts := range retention[0] {
					So(len(counts), ShouldEqual, 5)

					So(counts[0], ShouldEqual, 1)
					So(counts[1], ShouldEqual, 1)
					So(counts[2], ShouldEqual, 1)
					So(counts[3], ShouldEqual, 0)
					So(counts[4], ShouldEqual, 0)
				}

				for _, counts := range retention[1] {
					So(len(counts), ShouldEqual, 5)

					So(counts[0], ShouldEqual, 2)
					So(counts[1], ShouldEqual, .5)
					So(counts[2], ShouldEqual, 0)
					So(counts[3], ShouldEqual, 0)
					So(counts[4], ShouldEqual, 0)
				}

				for _, counts := range retention[2] {
					So(len(counts), ShouldEqual, 4)

					So(counts[0], ShouldEqual, 1)
					So(counts[1], ShouldEqual, 0)
					So(counts[2], ShouldEqual, 0)
					So(counts[3], ShouldEqual, 0)
				}

				for _, counts := range retention[3] {
					So(len(counts), ShouldEqual, 3)

					So(counts[0], ShouldEqual, 0)
					So(counts[1], ShouldEqual, 0)
					So(counts[2], ShouldEqual, 0)
				}

				for _, counts := range retention[4] {
					So(len(counts), ShouldEqual, 2)

					So(counts[0], ShouldEqual, 2)
					So(counts[1], ShouldEqual, 1)
				}

				for _, counts := range retention[5] {
					So(len(counts), ShouldEqual, 1)

					So(counts[0], ShouldEqual, 2)
				}

				Reset(func() { client.store.Do("FLUSHALL") })
			})
		})
	})
}
