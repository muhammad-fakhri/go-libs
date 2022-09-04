package cacheutils_test

import (
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-multierror"
	"github.com/muhammad-fakhri/go-libs/cache/mock_cache"
	"github.com/muhammad-fakhri/go-libs/cacheutils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHashWrapper(t *testing.T) {
	Convey("CachedHGet()", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHashCache := mock_cache.NewMockHashCacher(ctrl)
		mockHashCache.EXPECT().ErrorOnHashCacheMiss().Return(ErrCacheMiss).AnyTimes()

		advancedWrapper := cacheutils.NewHashWrapper(mockHashCache, mockHashCache, json.Marshal, json.Unmarshal)

		Convey("When it is a cache miss", func() {
			mockHashCache.EXPECT().HGet(gomock.Any(), gomock.Any()).Return("", ErrCacheMiss).Times(1)
			mockHashCache.EXPECT().HSet(gomock.Eq("key"), gomock.Eq("field"), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			fn := func() (interface{}, error) {
				return "data", nil
			}

			data, err := advancedWrapper.CachedHGet("key", "field", 0, fn, "")

			So(data, ShouldEqual, "data")
			So(err, ShouldBeNil)
		})

		Convey("When fetching from cache fails", func() {
			mockHashCache.EXPECT().HGet(gomock.Any(), gomock.Any()).Return("", ErrCacheFail).Times(1)
			mockHashCache.EXPECT().HSet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(0)
			fn := func() (interface{}, error) {
				return "data", nil
			}

			data, err := advancedWrapper.CachedHGet("key", "field", 0, fn, "")

			So(data, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		Convey("When it is a cache hit", func() {
			cached, _ := json.Marshal("cached")
			mockHashCache.EXPECT().HGet(gomock.Eq("key"), gomock.Eq("field")).Return(string(cached), nil).Times(1)
			mockHashCache.EXPECT().HSet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(0)
			fn := func() (interface{}, error) {
				return "data", nil
			}

			data, err := advancedWrapper.CachedHGet("key", "field", 0, fn, "")

			So(data, ShouldEqual, "cached")
			So(err, ShouldBeNil)
		})
	})

	Convey("Invalidate()", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCache := mock_cache.NewMockHashCacher(ctrl)
		advancedWrapper := cacheutils.NewHashWrapper(mockCache, mockCache, json.Marshal, json.Unmarshal)

		Convey("When cache deletion is successful", func() {
			mockCache.EXPECT().HDel(gomock.Any(), gomock.Any()).Return(int64(0), nil).Times(2)

			fields := []string{"a", "b"}
			err := advancedWrapper.Invalidate("key", fields...)

			So(err, ShouldBeNil)
		})

		Convey("When cache deletion is unsuccessful", func() {
			mockCache.EXPECT().HDel(gomock.Any(), gomock.Any()).Return(int64(0), ErrCacheFail).Times(2)

			fields := []string{"a", "b"}
			err := advancedWrapper.Invalidate("key", fields...)
			So(err, ShouldNotBeNil)

			merr, ok := err.(*multierror.Error)
			So(ok, ShouldBeTrue)
			So(merr.Errors, ShouldHaveLength, 2)
		})
	})
}
