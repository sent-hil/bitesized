package bitesized

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNearestInterval(t *testing.T) {
	tt := time.Date(1981, time.June, 12, 01, 1, 0, 0, time.UTC)

	Convey("It should find nearest hour", t, func() {
		n := nearestInterval(tt, Hour)
		So(n, ShouldEqual, "hour:1981-06-12 01:00")
	})

	Convey("It should find nearest day", t, func() {
		n := nearestInterval(tt, Day)
		So(n, ShouldEqual, "day:1981-06-12")
	})

	Convey("It should find nearest week", t, func() {
		n := nearestInterval(tt, Week)
		So(n, ShouldEqual, "week:1981-06-07")
	})

	Convey("It should find nearest month", t, func() {
		n := nearestInterval(tt, Month)
		So(n, ShouldEqual, "month:1981-06")
	})

	Convey("It should find nearest year", t, func() {
		n := nearestInterval(tt, Year)
		So(n, ShouldEqual, "year:1981")
	})
}
