package cache

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	convey.Convey("test new cache", t, func() {
		convey.Convey("test new redigo", func() {
			c, err := New(Redis, &Config{})
			convey.So(c, convey.ShouldNotBeNil)
			convey.So(err, convey.ShouldNotBeNil)
			c, err = New(123, &Config{})
			convey.So(c, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
		})
	})
}
