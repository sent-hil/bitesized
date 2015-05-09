package bitesized

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
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

		Reset(func() { client.store.Do("FLUSHALL") })
	})
}

func TestTrackEvent(t *testing.T) {
	Convey("It should return error unless event or username", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		err = client.TrackEvent("", "indianajones", time.Now())
		So(err, ShouldEqual, ErrInvalidArg)

		err = client.TrackEvent("dodge rock", "", time.Now())
		So(err, ShouldEqual, ErrInvalidArg)
	})

	Convey("It should track event for single interval", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		client.Intervals = []Interval{Year}

		err = client.TrackEvent("dodge rock", "indianajones", randomTime)
		So(err, ShouldBeNil)

		bitvalue, err := redis.Int(client.store.Do("GETBIT", "bitesized:year:1981", 1))
		So(err, ShouldBeNil)
		So(bitvalue, ShouldEqual, 1)

		Reset(func() { client.store.Do("FLUSHALL") })
	})

	Convey("It should track event for multiple intervals", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		client.Intervals = []Interval{Hour, Day, Week, Month, Year}

		err = client.TrackEvent("dodge rock", "indianajones", randomTime)
		So(err, ShouldBeNil)

		keys := []string{
			"bitesized:hour:1981-06-12-01:00",
			"bitesized:day:1981-06-12",
			"bitesized:week:1981-06-07",
			"bitesized:month:1981-06",
			"bitesized:year:1981",
		}

		for _, k := range keys {
			bitvalue, err := redis.Int(client.store.Do("GETBIT", k, 1))
			So(err, ShouldBeNil)
			So(bitvalue, ShouldEqual, 1)
		}

		Reset(func() { client.store.Do("FLUSHALL") })
	})
}

func TestKeyBuilder(t *testing.T) {
	Convey("It should return prefix with suffix if prefix", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		client.KeyPrefix = "prefix"
		So(client.key("suffix"), ShouldEqual, "prefix:suffix")

		Reset(func() { client.store.Do("FLUSHALL") })
	})

	Convey("It should return just suffix if no prefix", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		client.KeyPrefix = ""
		So(client.key("suffix"), ShouldEqual, "suffix")

		Reset(func() { client.store.Do("FLUSHALL") })
	})
}

func TestUserListKey(t *testing.T) {
	Convey("It should return user list key", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		client.KeyPrefix = ""
		So(client.userListKey(), ShouldEqual, "user-list")

		Reset(func() { client.store.Do("FLUSHALL") })
	})
}

func TestGetOrSetUser(t *testing.T) {
	username := "indianajones"

	Convey("It should save user if new user", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		id, err := client.getOrSetUser(username)
		So(err, ShouldBeNil)
		So(id, ShouldEqual, 1)

		Convey("It should get user if existing user", func() {
			id, err := client.getOrSetUser(username)
			So(err, ShouldBeNil)
			So(id, ShouldEqual, 1)
		})

		Reset(func() { client.store.Do("FLUSHALL") })
	})
}
