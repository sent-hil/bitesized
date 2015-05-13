package bitesized

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTrackEvent(t *testing.T) {
	Convey("It should return error unless event or user", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		err = client.TrackEvent("", user, time.Now())
		So(err, ShouldEqual, ErrInvalidArg)

		err = client.TrackEvent("dodge", "", time.Now())
		So(err, ShouldEqual, ErrInvalidArg)
	})

	Convey("It should track event for single interval", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		client.Intervals = []Interval{Year}

		err = client.TrackEvent("dodge rock", user, randomTime)
		So(err, ShouldBeNil)

		bitvalue, err := redis.Int(client.store.Do("GETBIT", "bitesized:event:dodge-rock:year:1981", 1))
		So(err, ShouldBeNil)
		So(bitvalue, ShouldEqual, 1)

		Reset(func() { client.store.Do("FLUSHALL") })
	})

	Convey("It should track event for multiple intervals", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		client.Intervals = []Interval{Hour, Day, Week, Month, Year}

		err = client.TrackEvent("dodge rock", user, randomTime)
		So(err, ShouldBeNil)

		keys := []string{
			"bitesized:event:dodge-rock:hour:1981-06-12-01:00",
			"bitesized:event:dodge-rock:day:1981-06-12",
			"bitesized:event:dodge-rock:week:1981-06-07",
			"bitesized:event:dodge-rock:month:1981-06",
			"bitesized:event:dodge-rock:year:1981",
		}

		for _, k := range keys {
			bitvalue, err := redis.Int(client.store.Do("GETBIT", k, 1))
			So(err, ShouldBeNil)
			So(bitvalue, ShouldEqual, 1)
		}

		Reset(func() { client.store.Do("FLUSHALL") })
	})
}

func TestCountEvent(t *testing.T) {
	Convey("", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		Convey("It should return 0 if no user did event", func() {
			count, err := client.CountEvent("dodge rock", time.Now(), Hour)
			So(err, ShouldBeNil)

			So(count, ShouldEqual, 0)
		})

		Convey("It should return count of users who did event", func() {
			err := client.TrackEvent("dodge rock", user, time.Now())
			So(err, ShouldBeNil)

			count, err := client.CountEvent("dodge rock", time.Now(), Hour)
			So(err, ShouldBeNil)

			So(count, ShouldEqual, 1)

			Reset(func() { client.store.Do("FLUSHALL") })
		})
	})
}

func TestDidEvent(t *testing.T) {
	Convey("", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		client.Intervals = []Interval{Hour}

		Convey("It should return no if user didn't do event", func() {
			didEvent, err := client.DidEvent("dodge rock", user, time.Now(), Hour)
			So(err, ShouldBeNil)

			So(didEvent, ShouldBeFalse)
		})

		Convey("It should return yes if user did event", func() {
			err = client.TrackEvent("dodge rock", user, time.Now())
			So(err, ShouldBeNil)

			didEvent, err := client.DidEvent("dodge rock", user, time.Now(), Hour)
			So(err, ShouldBeNil)

			So(didEvent, ShouldBeTrue)

			Reset(func() { client.store.Do("FLUSHALL") })
		})
	})
}

func TestGetEvents(t *testing.T) {
	Convey("", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		err = client.TrackEvent("dodge rock", user, time.Now())
		So(err, ShouldBeNil)

		err = client.TrackEvent("something other thing", user, time.Now())
		So(err, ShouldBeNil)

		Convey("It should return list of all events", func() {
			events, err := client.GetEvents("*")
			So(err, ShouldBeNil)

			So(len(events), ShouldEqual, 2)
		})

		Convey("It should return list of events with prefix", func() {
			events, err := client.GetEvents("dodge*")
			So(err, ShouldBeNil)

			So(len(events), ShouldEqual, 1)
		})

		Convey("It should return list of events when no prefix", func() {
			client, err := NewClient(testredis)
			So(err, ShouldBeNil)

			client.KeyPrefix = ""

			err = client.TrackEvent("dodge rock", user, time.Now())
			So(err, ShouldBeNil)

			events, err := client.GetEvents("dodge*")
			So(err, ShouldBeNil)

			So(len(events), ShouldEqual, 1)
		})

		Reset(func() { client.store.Do("FLUSHALL") })
	})
}
