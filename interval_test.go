package bitesized

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var randomTime = time.Date(1981, time.June, 12, 01, 42, 0, 0, time.UTC)

func TestNearestInterval(t *testing.T) {
	Convey("It should empty for 'All'", t, func() {
		n := nearestInterval(randomTime, All)
		So(n, ShouldEqual, "all")
	})

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

	Convey("It should find nearest quarter", t, func() {
		n := nearestInterval(randomTime, Quarter)
		So(n, ShouldEqual, "quarter:1981-04")
	})

	Convey("It should find nearest 10 minute cycle", t, func() {
		n := nearestInterval(randomTime, TenMinutes)
		So(n, ShouldEqual, "ten_minutes:1981-06-12-01:40")
	})

	Convey("It should find nearest 30 minute cycle", t, func() {
		n := nearestInterval(randomTime, ThirtyMinutes)
		So(n, ShouldEqual, "thirty_minutes:1981-06-12-01:30")
	})

	Convey("It should find nearest biweekly date (first part)", t, func() {
		testingTime := time.Date(1981, time.June, 12, 01, 42, 0, 0, time.UTC)
		n := nearestInterval(testingTime, Biweekly)
		So(n, ShouldEqual, "biweekly:1981-06-10")
	})

	Convey("It should find nearest biweekly date (second part)", t, func() {
		testingTime := time.Date(1981, time.June, 16, 01, 42, 0, 0, time.UTC)
		n := nearestInterval(testingTime, Biweekly)
		So(n, ShouldEqual, "biweekly:1981-06-14")
	})

	Convey("It should find nearest bimonthly date (first part)", t, func() {
		testingTime := time.Date(1981, time.June, 12, 01, 42, 0, 0, time.UTC)
		n := nearestInterval(testingTime, Bimonthly)
		So(n, ShouldEqual, "bimonthly:1981-06-01")
	})

	Convey("It should find nearest bimonthly date (second part)", t, func() {
		testingTime := time.Date(1981, time.June, 28, 01, 42, 0, 0, time.UTC)
		n := nearestInterval(testingTime, Bimonthly)
		So(n, ShouldEqual, "bimonthly:1981-06-15")
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
