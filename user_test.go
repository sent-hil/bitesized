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

func TestEventUsers(t *testing.T) {
	Convey("It should return list of users who did an event", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		client.Intervals = []Interval{Hour}

		err = client.TrackEvent("dodge rock", user, randomTime)
		So(err, ShouldBeNil)

		err = client.TrackEvent("dodge rock", user+"1", randomTime)
		So(err, ShouldBeNil)

		users, err := client.EventUsers("dodge rock", randomTime, Hour)
		So(err, ShouldBeNil)

		So(len(users), ShouldEqual, 2)
		So(users[0], ShouldEqual, user)
		So(users[1], ShouldEqual, user+"1")

		Reset(func() { client.store.Do("FLUSHALL") })
	})
}

func TestRemoveUser(t *testing.T) {
	Convey("It should remove user", t, func() {
		client, err := NewClient(testredis)
		So(err, ShouldBeNil)

		client.Intervals = []Interval{Hour, Day}

		err = client.TrackEvent("dodge rock", user, randomTime)
		So(err, ShouldBeNil)

		err = client.RemoveUser(user)
		So(err, ShouldBeNil)

		didEvent, err := client.DidEvent("dodge rock", user, randomTime, Hour)
		So(err, ShouldBeNil)

		So(didEvent, ShouldBeFalse)

		Reset(func() { client.store.Do("FLUSHALL") })
	})
}
