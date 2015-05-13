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
