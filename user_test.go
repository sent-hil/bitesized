package bitesized

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsUserNew(t *testing.T) {
	Convey("", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		Convey("It should return true if user is new", func() {
			isNew, err := client.IsUserNew(user)
			So(err, ShouldBeNil)

			So(isNew, ShouldBeTrue)
		})

		Convey("It should return false if user isn't new", func() {
			_, err := client.getOrSetUser(user)
			So(err, ShouldBeNil)

			isNew, err := client.IsUserNew(user)
			So(err, ShouldBeNil)

			So(isNew, ShouldBeFalse)

			Reset(func() { client.store.Do("FLUSHALL") })
		})
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

		Convey("It should get existing user by id", func() {
			id, err := client.getOrSetUser(user)
			So(err, ShouldBeNil)

			username, err := client.getUserById(id)
			So(err, ShouldBeNil)
			So(username, ShouldEqual, user)
		})

		Reset(func() { client.store.Do("FLUSHALL") })
	})
}
