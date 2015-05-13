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

func TestOperation(t *testing.T) {
	Convey("", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		Convey("It should do specified operation", func() {
			k1 := "testkey1"
			_, err = client.store.Do("SETBIT", k1, 1, On)
			So(err, ShouldBeNil)

			k2 := "testkey2"
			_, err = client.store.Do("SETBIT", k2, 2, On)
			So(err, ShouldBeNil)

			keys := []string{k1, k2}

			count, err := client.Operation(AND, keys...)
			So(err, ShouldBeNil)
			So(count, ShouldEqual, 0)

			count, err = client.Operation(OR, keys...)
			So(err, ShouldBeNil)
			So(count, ShouldEqual, 2)

			count, err = client.Operation(XOR, keys...)
			So(err, ShouldBeNil)
			So(count, ShouldEqual, 2)

			count, err = client.Operation(NOT, k1)
			So(err, ShouldBeNil)
			So(count, ShouldEqual, 7)

			Reset(func() { client.store.Do("FLUSHALL") })
		})

		Convey("It should accept only one op for NOT", func() {
			_, err := client.Operation(NOT, "k1", "k2")
			So(err, ShouldEqual, ErrNotOpAcceptsOnekey)
		})
	})
}
