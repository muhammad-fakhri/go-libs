package cache

//go:generate mockgen -destination mock_cache/mock_cache.go . Cache,Cacher,Conn,HashCacher,MultiCacher,Scripter

import (
	"errors"
	"time"

	goredis "github.com/go-redis/redis/v7"
)

// Cacher is an interface for basic caching operations.
type Cacher interface {
	GetConn() Conn
	Set(key, value string, ttl time.Duration) error
	Get(key string) (string, error)
	Del(key ...string) error
	ErrorOnCacheMiss() error
}

type Conn interface {
	Close() error

	Err() error

	Do(commandName string, args ...interface{}) (reply interface{}, err error)

	Send(commandName string, args ...interface{}) error

	Flush() error

	Receive() (reply interface{}, err error)
}

// HashCacher is an interface for redis hash operations
type HashCacher interface {
	// HSet Sets field in the hash stored at key to value. If key does not exist, a new key holding a hash is created. If field already exists in the hash, it is overwritten.
	HSet(key, field, value string, ttl time.Duration) error
	// HSetNX Sets field in the hash stored at key to value, only if field does not yet exist.
	HSetNX(key, field, value string, ttl time.Duration) error
	// HMSet Sets multiple pair of field & value
	HMSet(key string, fieldsMap map[string]string, ttl time.Duration) error
	// HGet Returns the value associated with field in the hash stored at key.
	HGet(key, field string) (string, error)
	// HMGet Returns values with multiple fields
	HMGet(key string, fields ...string) ([]string, error)
	// HDel Removes the specified fields from the hash stored at key. Specified fields that do not exist within this hash are ignored. If key does not exist, it is treated as an empty hash and this command returns 0.
	HDel(key string, fields ...string) (int64, error)
	// HKeys Returns all field names in the hash stored at key.
	HKeys(key string) ([]string, error)
	// HVals Returns all values in the hash stored at key.
	HVals(key string) ([]string, error)
	// HGetAll Returns all fields and values of the hash stored at key. In the returned value, every field name is followed by its value, so the length of the reply is twice the size of the hash.
	HGetAll(key string) (map[string]string, error)
	HExists(key, field string) (bool, error)
	HIncrBy(key, field string, incrValue int64) (int64, error)
	ErrorOnHashCacheMiss() error
}

// MultiCacher do multi cache operation (get and set multiple key in one opeartion).
type MultiCacher interface {
	// MSet sets multiple key into cache. input is a map of key:value
	MSet(values map[string]string) error
	// MSetEx sets multiple key into cache and set its expiry
	MSetEx(values map[string]string, ttl time.Duration) error
	// MGet gets multiple key from cache
	MGet(keys []string) ([]string, error)
}

// TODO: Should probably rename this to RedisClienter?
// This looks more like a redis client interface rather than a generic cache interface, ex: memcached wouldn't be able to implement all of this.
type Cache interface {
	Cacher
	HashCacher
	Scripter
	MultiCacher
	// SetNX et key to hold string value if key does not exist
	SetNX(key, value string, ttl time.Duration) error
	// ScanKeys get all key that match pattern
	ScanKeys(pattern string) ([]string, error)
	// IncrBy increments the number stored at key by increment. If the key does not exist, it is set to 0 before performing the operation
	IncrBy(key string, incr int64) (int64, error)
	// ZAdd add member to sorted set.
	ZAdd(key, member string, score int) error
	// ZAddXX add member to sorted set. Score is updated with new value. Never add score.
	ZAddXX(key, member string, score int) error
	// ZAddNX add member to sorted set. Score is added with new value. Never update score.
	ZAddNX(key, member string, score int64) (int64, error)
	// ZAddINCR add member to sorted set. Score is incremented from previous value if exist
	ZAddINCR(key, member string, score int) error
	// ZRange return members with its score from sorted set. Sorted ascending
	ZCard(key string) (int64, error)
	// ZCard return cardinality (number of elements) of sorted set.
	ZRange(key string, start, stop int) ([]string, error)
	// ZRevRange return members with its score from sorted set. Sorted descending
	ZRevRange(key string, start, stop int) ([]string, error)
	// ZRangeByScore return members with its score from sorted set, using score as range. Sorted ascending.
	ZRangeByScore(key string, min, max, offset, count int) ([]string, error)
	// ZRevRangeByScore return members with its score from sorted set, using score as range. Sorted ascending.
	ZRevRangeByScore(key string, max, min, offset, count int) ([]string, error)
	// ZRank return member rank in sorted set. Sorted ascending.
	ZRank(key, member string) (int64, error)
	// ZRevRank return member rank in sorted set. Sorted descending.
	ZRevRank(key, member string) (int64, error)
	// ZScore return members score
	ZScore(key, member string) (int64, error)
	// ZCount return count of members which score is between min and max.
	ZCount(key string, min, max int) (int64, error)
	// ZRemRangeByScore Delete member by score range
	ZRemRangeByScore(key string, start, stop int) (int64, error)
	// SAdd Add the specified members to the set stored at key. Specified members that are already a member of this set are ignored. If key does not exist, a new set is created before adding the specified members
	SAdd(key, member string) (int64, error)
	// SCard Returns the set cardinality (number of elements) of the set stored at key
	SCard(key string) (int64, error)
	// SDiff Returns the members of the set resulting from the difference between the first set and all the successive sets.
	SDiff(keys ...string) ([]string, error)
	// SDiff This command is equal to SDIFF, but instead of returning the resulting set, it is stored in destination.
	SDiffStore(keys ...string) (int64, error)
	// SInter Returns the members of the set resulting from the intersection of all the given sets.
	SInter(keys ...string) ([]string, error)
	// SInterStore This command is equal to SINTER, but instead of returning the resulting set, it is stored in destination.
	SInterStore(keys ...string) (int64, error)
	// SIsMember Returns if member is a member of the set stored at key.
	SIsMember(keys, member string) (int64, error)
	// SMembers Returns all the members of the set value stored at key. This has the same effect as running SINTER with one argument key.
	SMembers(key string) ([]string, error)
	// SMove Move member from the set at source to the set at destination. This operation is atomic. In every given moment the element will appear to be a member of source or destination for other clients.
	SMove(value, source, destination string) (int64, error)
	// SPop Removes and returns one or more random elements from the set value store at key.
	SPop(key string, count int) ([]string, error)
	// SRandMember When called with just the key argument, return a random element from the set value stored at key.
	SRandMember(key string, count int) ([]string, error)
	// SRem Remove the specified members from the set stored at key. Specified members that are not a member of this set are ignored
	SRem(key string, member string) (int64, error)
	// SUnion Returns the members of the set resulting from the union of all the given sets
	SUnion(keys ...string) ([]string, error)
	// This command is equal to SUNION, but instead of returning the resulting set, it is stored in destination.
	SUnionStore(keys ...string) (int64, error)
	ZRem(key string, members ...string) (int64, error)
	ZAddXXIncrBy(key, member string, incrValue int64) (int64, error)
	Expire(key string, ttl time.Duration) (int64, error)
	Exists(key string) (bool, error)
	ZRevRangeWithScore(key string, start, stop int64) (interface{}, error)
	// GeoAdd Adds the specified geospatial items to the specified key
	GeoAdd(key string, geos ...*GeoPoint) (int64, error)
	//GeoHash return valid Geohash strings representing the position of one or more elements in a sorted set value representing a geospatial index
	GeoHash(key string, members ...string) ([]string, error)
	// GeoRadius get items from given key within given radius from particular given point
	GeoRadius(key string, long, lat float64, q *GeoRadiusQuery) ([]*GeoLoc, error)
	// TTL returns remaining time to live of a key that has a timeout, return -1 if key doesn't have ttl
	TTL(key string) (int64, error)
	LLen(key string) (int64, error)
	LPop(key string, count int) ([]string, error)
	LPush(key string, values []string) (int64, error)
	LPushX(key string, values []string) (int64, error)
	RPop(key string, count int) ([]string, error)
	RPush(key string, values []string) (int64, error)
	RPushX(key string, values []string) (int64, error)
}

// Scripter is interface contract for redis scripting
type Scripter interface {
	IncrXX(key string, value int64) (reply int64, err error)
	DecrWithLimit(key string, value, lowerBound int64) (reply int64, err error)
	HGetSet(key, field, value, prevValue string, ttl time.Duration) error
	// ZAddToFixed add to sorted set with defined max size. If sorted set size reach max size, remove member with lowest score
	ZAddToFixed(key, member string, score, maxSize int) (reply int64, err error)
}

type Implementation int

const (
	Redis = Implementation(iota)
	RedisCluster
	RedisSentinel
)

type Config struct {
	ServerAddr      string
	MaxIdle         int
	MaxActive       int
	IdleTimeout     time.Duration
	MaxConnLifetime time.Duration
	Wait            bool
	UseCommonErr    bool
}

type ConfigCluster struct {
	StartNodes   []string
	KeepAlive    int
	ConnTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	AliveTime    time.Duration
}

// https://godoc.org/github.com/go-redis/redis#ClusterOptions
type ConfigCacheCluster struct {
	Addrs          []string
	MaxRedirects   int
	ReadOnly       bool
	RouteByLatency bool
	RouteRandomly  bool
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	PoolSize       int
	DialTimeout    time.Duration
	IdleTimeout    time.Duration
}

type FailoverOptions goredis.FailoverOptions

var (
	ErrNil                  = errors.New("redis status: nil")
	ErrInsufficientArgument = errors.New("wrong number of arguments")
)

// New return ready to use Cache instance
func New(impl Implementation, cfg interface{}) (Cache, error) {
	switch impl {
	case Redis:
		return newRedigo(cfg.(*Config))
	case RedisCluster:
		return newRedisCluster(cfg.(*ConfigCacheCluster))
	case RedisSentinel:
		return newRedisSentinel(cfg.(*FailoverOptions))
	}

	return nil, errors.New("no cache implementations found")
}
