package bitesized

import (
	"testing"

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

		Convey("It should save user if new user", func() {
			id, err = client.getOrSetUser(user + "1")
			So(err, ShouldBeNil)
			So(id, ShouldEqual, 2)
		})

		Convey("It should get user if existing user", func() {
			id, err := client.getOrSetUser(user)
			So(err, ShouldBeNil)
			So(id, ShouldEqual, 1)
		})

		Reset(func() { client.store.Do("FLUSHALL") })
	})
}
