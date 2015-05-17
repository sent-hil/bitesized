package bitesized

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var randomTime = time.Date(1981, time.June, 12, 01, 1, 0, 0, time.UTC)

func TestNearestInterval(t *testing.T) {
	Convey("It should find nearest hour", t, func() {
		n := nearestInterval(randomTime, Hour)
		So(n, ShouldEqual, "hour:1981-06-12-01:00")
	})

	Convey("It should find nearest day", t, func() {
		n := nearestInterval(randomTime, Day)
		So(n, ShouldEqual, "day:1981-06-12")
	})

	Convey("It should find nearest week", t, func() {
		n := nearestInterval(randomTime, Week)
		So(n, ShouldEqual, "week:1981-06-07")
	})

	Convey("It should find nearest month", t, func() {
		n := nearestInterval(randomTime, Month)
		So(n, ShouldEqual, "month:1981-06")
	})

	Convey("It should find nearest year", t, func() {
		n := nearestInterval(randomTime, Year)
		So(n, ShouldEqual, "year:1981")
	})
}

func TestGetDuration(t *testing.T) {
	Convey("It should return duration for hour", t, func() {
		d := getDuration(randomTime, Hour)
		So(d, ShouldEqual, 1*time.Hour)
	})

	Convey("It should return duration for day", t, func() {
		d := getDuration(randomTime, Day)
		So(d, ShouldEqual, 24*time.Hour)
	})

	Convey("It should return duration for week", t, func() {
		d := getDuration(randomTime, Week)
		So(d, ShouldEqual, 7*24*time.Hour)
	})

	Convey("It should return duration for month with 31 days", t, func() {
		t := time.Date(2015, time.January, 01, 00, 0, 0, 0, time.UTC)
		d := getDuration(t, Month)

		So(d, ShouldEqual, 31*24*time.Hour)
	})

	Convey("It should return duration for month with 30 days", t, func() {
		t := time.Date(2015, time.April, 01, 00, 0, 0, 0, time.UTC)
		d := getDuration(t, Month)

		So(d, ShouldEqual, 30*24*time.Hour)
	})

	Convey("It should return duration for month with 28 days", t, func() {
		t := time.Date(2015, time.February, 01, 00, 0, 0, 0, time.UTC)
		d := getDuration(t, Month)

		So(d, ShouldEqual, 28*24*time.Hour)
	})
}
