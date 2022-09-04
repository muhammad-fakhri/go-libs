package cacheutils_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/muhammad-fakhri/go-libs/cacheutils"
	"github.com/muhammad-fakhri/go-libscache/mock_cache"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-multierror"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	ErrCacheMiss = errors.New("cache miss")
	ErrCacheFail = errors.New("cache failure")
)

func TestWrapper(t *testing.T) {
	Convey("CachedGet()", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCache := mock_cache.NewMockCacher(ctrl)
		mockCache.EXPECT().ErrorOnCacheMiss().Return(ErrCacheMiss).AnyTimes()

		wrapper := cacheutils.NewWrapper(mockCache)

		Convey("When it is a cache miss", func() {
			mockCache.EXPECT().Get(gomock.Any()).Return("", ErrCacheMiss).Times(1)
			mockCache.EXPECT().Set(gomock.Eq("key"), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			fn := func() (interface{}, error) {
				return "data", nil
			}

			data, err := wrapper.CachedGet("key", 0, fn, "")

			So(data, ShouldEqual, "data")
			So(err, ShouldBeNil)
		})

		Convey("When fetching from cache fails", func() {
			mockCache.EXPECT().Get(gomock.Any()).Return("", ErrCacheFail).Times(1)
			mockCache.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(0)
			fn := func() (interface{}, error) {
				return "data", nil
			}

			data, err := wrapper.CachedGet("key", 0, fn, "")

			So(data, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		Convey("When it is a cache hit", func() {
			cached, _ := json.Marshal("cached")
			mockCache.EXPECT().Get(gomock.Eq("key")).Return(string(cached), nil).Times(1)
			mockCache.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(0)
			fn := func() (interface{}, error) {
				return "data", nil
			}

			data, err := wrapper.CachedGet("key", 0, fn, "")

			So(data, ShouldEqual, "cached")
			So(err, ShouldBeNil)
		})
	})

	Convey("Invalidate()", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCache := mock_cache.NewMockCacher(ctrl)
		wrapper := cacheutils.NewWrapper(mockCache)

		Convey("When cache deletion is successful", func() {
			mockCache.EXPECT().Del(gomock.Any()).Return(nil).Times(2)

			keys := []string{"a", "b"}
			err := wrapper.Invalidate(keys...)

			So(err, ShouldBeNil)
		})

		Convey("When cache deletion is unsuccessful", func() {
			mockCache.EXPECT().Del(gomock.Any()).Return(ErrCacheFail).Times(2)

			keys := []string{"a", "b"}
			err := wrapper.Invalidate(keys...)
			So(err, ShouldNotBeNil)

			merr, ok := err.(*multierror.Error)
			So(ok, ShouldBeTrue)
			So(merr.Errors, ShouldHaveLength, 2)
		})
	})
}
