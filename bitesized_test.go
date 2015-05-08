package bitesized

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var testredis = "localhost:6379"

func TestNewClient(t *testing.T) {
	Convey("It should initialize client", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		Convey("With redis connection", func() {
			So(client.store, ShouldNotBeNil)
		})

		Convey("With default values", func() {
			So(len(client.Intervals), ShouldBeGreaterThan, 1)
			So(client.KeyPrefix, ShouldEqual, "bitesized")
		})
	})
}

func TestTrackEvent(t *testing.T) {
	Convey("It should track event in single interval", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		client.Intervals = []Interval{Hour}

		err = client.TrackEvent("dodge rock", "indianajones", time.Now())
		So(err, ShouldBeNil)
	})
}
