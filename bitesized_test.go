package bitesized

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	testredis = "localhost:6379"
	user      = "indianajones"
)

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

		bitvalue, err := redis.Int(client.store.Do("GETBIT", "bitesized:dodge-rock:year:1981", 1))
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
			"bitesized:dodge-rock:hour:1981-06-12-01:00",
			"bitesized:dodge-rock:day:1981-06-12",
			"bitesized:dodge-rock:week:1981-06-07",
			"bitesized:dodge-rock:month:1981-06",
			"bitesized:dodge-rock:year:1981",
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

		Convey("It should join multiple suffixes", func() {
			So(client.key("one", "two"), ShouldEqual, "prefix:one:two")
		})
	})

	Convey("It should return just suffix if no prefix", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		client.KeyPrefix = ""
		So(client.key("suffix"), ShouldEqual, "suffix")

		Convey("It should join multiple suffixes", func() {
			So(client.key("one", "two"), ShouldEqual, "one:two")
		})
	})
}

func TestUserListKey(t *testing.T) {
	Convey("It should return user list key", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		client.KeyPrefix = ""
		So(client.userListKey(), ShouldEqual, "user-list")
	})
}

func TestGetOrSetUser(t *testing.T) {
	Convey("It should save user if new user", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		id, err := client.getOrSetUser(user)
		So(err, ShouldBeNil)
		So(id, ShouldEqual, 1)

		Convey("It should get user if existing user", func() {
			id, err := client.getOrSetUser(user)
			So(err, ShouldBeNil)
			So(id, ShouldEqual, 1)
		})

		Reset(func() { client.store.Do("FLUSHALL") })
	})
}

func TestCountEvent(t *testing.T) {
	Convey("It should return 0 if no user did event", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		count, err := client.CountEvent("dodge rock", time.Now(), Hour)
		So(err, ShouldBeNil)

		So(count, ShouldEqual, 0)
	})

	Convey("It should return count of users who did event", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		client.Intervals = []Interval{Hour}

		client.TrackEvent("dodge rock", user, time.Now())

		count, err := client.CountEvent("dodge rock", time.Now(), Hour)
		So(err, ShouldBeNil)

		So(count, ShouldEqual, 1)

		Reset(func() { client.store.Do("FLUSHALL") })
	})
}

func TestDasherize(t *testing.T) {
	Convey("It should split event on space and join with dash", t, func() {
		So(dasherize("dodge"), ShouldEqual, "dodge")
		So(dasherize("dodge rock"), ShouldEqual, "dodge-rock")
	})
}
