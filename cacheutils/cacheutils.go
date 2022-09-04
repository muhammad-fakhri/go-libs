package cacheutils

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/muhammad-fakhri/go-libs/cache"
)

type Wrapper interface {
	CachedGet(key string, ttl time.Duration, fn func() (interface{}, error), bufferForType interface{}) (interface{}, error)
	Invalidate(keys ...string) error
}

func NewWrapper(cache cache.Cacher) Wrapper {
	return &wrapper{cache: cache}
}

type wrapper struct {
	cache cache.Cacher
}

// Parameter bufferForType is used for type inference in json.Unmarshal()
func (w *wrapper) CachedGet(key string, ttl time.Duration, fn func() (interface{}, error), bufferForType interface{}) (interface{}, error) {
	var data interface{}
	cached, err := w.cache.Get(key)

	if err == w.cache.ErrorOnCacheMiss() { // cache miss
		data, err := fn()
		if err != nil {
			return nil, err
		}

		marshaled, err := json.Marshal(data)
		if err != nil {
			log.Println("Error marshaling data: ", data)
		}
		err = w.cache.Set(key, string(marshaled), ttl)
		if err != nil {
			log.Println("Unable to set cache with key: ", key)
		}

		return data, nil
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(cached), &bufferForType)
	if err != nil {
		return nil, err
	}

	data = bufferForType
	return data, nil
}

func (w *wrapper) Invalidate(keys ...string) error {
	var errors *multierror.Error

	for _, key := range keys {
		if err := w.cache.Del(key); err != nil {
			errors = multierror.Append(errors, err)
		}
	}

	return errors.ErrorOrNil()
}
