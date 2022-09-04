package cacheutils

import (
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/muhammad-fakhri/go-libs/cache"
)

type HashWrapper interface {
	CachedHGet(key, field string, ttl time.Duration, fn func() (interface{}, error), bufferForType interface{}) (interface{}, error)
	Invalidate(key string, fields ...string) error
}

func NewHashWrapper(cacheMaster, cacheSlave cache.HashCacher, marshal func(v interface{}) ([]byte, error), unmarshal func(data []byte, v interface{}) error) HashWrapper {
	return &hashWrapper{
		cacheMaster: cacheMaster,
		cacheSlave:  cacheSlave,
		marshal:     marshal,
		unmarshal:   unmarshal,
	}
}

type hashWrapper struct {
	cacheMaster cache.HashCacher
	cacheSlave  cache.HashCacher
	marshal     func(v interface{}) ([]byte, error)
	unmarshal   func(data []byte, v interface{}) error
}

// Parameter bufferForType is used for type inference in unmarshal()
func (w *hashWrapper) CachedHGet(key, field string, ttl time.Duration, fn func() (interface{}, error), bufferForType interface{}) (interface{}, error) {
	var data interface{}
	cached, err := w.cacheSlave.HGet(key, field)

	if err == w.cacheSlave.ErrorOnHashCacheMiss() { // cache miss
		data, err := fn()
		if err != nil {
			return nil, err
		}

		marshaled, err := w.marshal(data)
		if err != nil {
			log.Println("[Common][CachedHGet] Error marshaling data: ", data)
		}
		err = w.cacheMaster.HSet(key, field, string(marshaled), ttl)
		if err != nil {
			log.Printf("[Common][CachedHGet] Unable to HSET cache with key %s and field %s \n", key, field)
		}

		return data, nil
	} else if err != nil {
		return nil, err
	}

	err = w.unmarshal([]byte(cached), &bufferForType)
	if err != nil {
		return nil, err
	}

	data = bufferForType
	return data, nil
}

func (w *hashWrapper) Invalidate(key string, fields ...string) error {
	var errors *multierror.Error

	for _, field := range fields {
		if _, err := w.cacheMaster.HDel(key, field); err != nil {
			errors = multierror.Append(errors, err)
		}
	}

	return errors.ErrorOrNil()
}
