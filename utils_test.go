package bitesized

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDasherize(t *testing.T) {
	Convey("It should split event on space and join with dash", t, func() {
		So(dasherize("dodge"), ShouldEqual, "dodge")
		So(dasherize("dodge rock"), ShouldEqual, "dodge-rock")
	})
}

func TestRandomSeq(t *testing.T) {
	Convey("It should return random string", t, func() {
		So(randSeq(20), ShouldNotEqual, randSeq(2))
	})
}
