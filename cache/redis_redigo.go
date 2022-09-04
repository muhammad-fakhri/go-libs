package cache

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	rediscluster "github.com/marioorlando/redis-go-cluster"
)

type redigoImpl struct {
	Pool         *redis.Pool
	Cluster      *rediscluster.Cluster
	UseCommonErr bool
}

const (
	redisSet              = "SET"
	redisGet              = "GET"
	redisMSet             = "MSET"
	redisMGet             = "MGET"
	redisHGet             = "HGET"
	redisHMGet            = "HMGET"
	redisHSet             = "HSET"
	redisHSetNX           = "HSETNX"
	redisHMSet            = "HMSET"
	redisHKeys            = "HKEYS"
	redisHVals            = "HVALS"
	redisHGetAll          = "HGETALL"
	redisHLen             = "HLEN"
	redisScan             = "SCAN"
	redisMatch            = "MATCH"
	redisIncrBy           = "INCRBY"
	redisDel              = "DEL"
	redisHDel             = "HDEL"
	redisEx               = "EX"
	redisNX               = "NX"
	redisXX               = "XX"
	redisLimit            = "LIMIT"
	redisPing             = "PING"
	redisExpire           = "EXPIRE"
	redisTTL              = "TTL"
	redisINCR             = "INCR"
	redisZAdd             = "ZADD"
	redisZCard            = "ZCARD"
	redisZRange           = "ZRANGE"
	redisZRevRange        = "ZREVRANGE"
	redisZRangeByScore    = "ZRANGEBYSCORE"
	redisZRevRangeByScore = "ZREVRANGEBYSCORE"
	redisZRank            = "ZRANK"
	redisZRevRank         = "ZREVRANK"
	redisZScore           = "ZSCORE"
	redisZCount           = "ZCOUNT"
	redisZRemRangeByScore = "ZREMRANGEBYSCORE"
	redisSAdd             = "SADD"
	redisSCard            = "SCARD"
	redisSDiff            = "SDIFF"
	redisSDiffStore       = "SDIFFSTORE"
	redisSInter           = "SINTER"
	redisSInterStore      = "SINTERSTORE"
	redisSIsMember        = "SISMEMBER"
	redisSMembers         = "SMEMBERS"
	redisSMove            = "SMOVE"
	redisSPop             = "SPOP"
	redisSRandMember      = "SRANDMEMBER"
	redisSRem             = "SREM"
	redisSUnion           = "SUNION"
	redisSUnionStore      = "SUNIONSTORE"
	redisZREM             = "ZREM"
	redisWithScores       = "WITHSCORES"
	redisHExists          = "HEXISTS"
	redisHIncrBy          = "HINCRBY"
	redisExists           = "EXISTS"
	redisGeoAdd           = "GEOADD"
	redisGeoHash          = "GEOHASH"
	redisGeoRadius        = "GEORADIUS"
	redisMulti            = "MULTI"
	redisExec             = "EXEC"
	redisLLen             = "LLEN"
	redisLPop             = "LPOP"
	redisLPush            = "LPUSH"
	redisLPushX           = "LPUSHX"
	redisRPop             = "RPOP"
	redisRPush            = "RPUSH"
	redisRPushX           = "RPUSHX"
)

var (
	ErrNX  = errors.New("key already exist")
	ErrXX  = errors.New("key is not exists")
	ErrHNX = errors.New("(key,field) combination already exists")
)

func newRedigo(cfg *Config) (*redigoImpl, error) {
	r := &redigoImpl{
		Pool: &redis.Pool{
			IdleTimeout:     cfg.IdleTimeout,
			MaxActive:       cfg.MaxActive,
			MaxConnLifetime: cfg.MaxConnLifetime,
			MaxIdle:         cfg.MaxIdle,
			Wait:            cfg.Wait,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", cfg.ServerAddr)
				if err != nil {
					return nil, err
				}

				return c, nil
			},
		},
		UseCommonErr: cfg.UseCommonErr,
	}
	_, err := r.Pool.Get().Do(redisPing)

	return r, err
}

func newRedigoCluster(cfg *ConfigCluster) (*redigoImpl, error) {
	cluster, err := rediscluster.NewCluster(&rediscluster.Options{
		StartNodes:   cfg.StartNodes,
		ConnTimeout:  cfg.ConnTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		KeepAlive:    cfg.KeepAlive,
		AliveTime:    cfg.AliveTime,
	})
	if err != nil {
		return nil, err
	}
	r := &redigoImpl{
		Cluster: cluster,
	}
	return r, nil
}

func (r *redigoImpl) ErrorOnHashCacheMiss() error {
	if r.UseCommonErr {
		return ErrNil
	}
	return redis.ErrNil
}

func (r *redigoImpl) ErrorOnCacheMiss() error {
	if r.UseCommonErr {
		return ErrNil
	}
	return redis.ErrNil
}

func (r *redigoImpl) GetConn() Conn {
	return r.Pool.Get()
}

func (r *redigoImpl) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if r.Pool != nil {
		c := r.GetConn()
		defer c.Close()
		reply, err = c.Do(commandName, args...)
	} else {
		reply, err = r.Cluster.Do(commandName, args...)
	}
	return
}

func (r *redigoImpl) Set(key, value string, ttl time.Duration) error {
	var err error

	if 0 >= ttl {
		_, err = r.Do(redisSet, key, value)
	} else {
		_, err = r.Do(redisSet, key, value, redisEx, int64(ttl.Seconds()))
	}

	return err
}

func (r *redigoImpl) SetNX(key, value string, ttl time.Duration) error {
	var (
		err   error
		reply interface{}
	)

	if 0 >= ttl {
		reply, err = r.Do(redisSet, key, value, redisNX)
	} else {
		reply, err = r.Do(redisSet, key, value, redisEx, int64(ttl.Seconds()), redisNX)
	}

	if nil != err {
		return err
	}
	if nil == reply {
		return ErrNX
	}

	return err
}

func (r *redigoImpl) ScanKeys(pattern string) ([]string, error) {
	iter := 0
	keys := []string{}
	for {
		arr, err := redis.Values(r.Do(redisScan, iter, redisMatch, pattern))
		if err != nil {
			return nil, err
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}
	return keys, nil
}

func (r *redigoImpl) HSet(key, field, value string, ttl time.Duration) error {
	var err error

	_, err = r.Do(redisHSet, key, field, value)
	if err != nil {
		return err
	}
	if ttl > 0 {
		_, err = r.Do(redisExpire, key, int64(ttl.Seconds()))
	}

	return err
}

func (r *redigoImpl) HSetNX(key, field, value string, ttl time.Duration) error {
	var err error

	reply, err := redis.Int64(r.Do(redisHSetNX, key, field, value))
	if err != nil {
		return err
	}
	if reply == 0 {
		return ErrHNX
	}

	if ttl > 0 {
		_, err = r.Do(redisExpire, key, int64(ttl.Seconds()))
	}

	return err
}

func (r *redigoImpl) HMSet(key string, fieldsMap map[string]string, ttl time.Duration) error {
	_, err := r.Do(redisHMSet, redis.Args{}.Add(key).AddFlat(fieldsMap)...)
	if err != nil {
		return err
	}

	if ttl > 0 {
		_, err = r.Do(redisExpire, key, int64(ttl.Seconds()))
	}

	return err
}

func (r *redigoImpl) HGet(key, field string) (string, error) {
	reply, err := redis.String(r.Do(redisHGet, key, field))
	if err == redis.ErrNil && r.UseCommonErr {
		return reply, ErrNil
	}
	return reply, err
}

func (r *redigoImpl) HMGet(key string, fields ...string) ([]string, error) {
	reply, err := r.Do(redisHMGet, redis.Args{}.Add(key).AddFlat(fields)...)
	return redis.Strings(reply, err)
}

func (r *redigoImpl) HDel(key string, fields ...string) (int64, error) {
	return redis.Int64(r.Do(redisHDel, redis.Args{}.Add(key).AddFlat(fields)...))
}

func (r *redigoImpl) HKeys(key string) ([]string, error) {
	reply, err := r.Do(redisHKeys, key)
	return redis.Strings(reply, err)
}

func (r *redigoImpl) HVals(key string) ([]string, error) {
	reply, err := r.Do(redisHVals, key)
	return redis.Strings(reply, err)
}

func (r *redigoImpl) HGetAll(key string) (map[string]string, error) {
	reply, err := r.Do(redisHGetAll, key)
	return redis.StringMap(reply, err)
}

func (r *redigoImpl) HLen(key, field string) (int64, error) {
	var (
		reply interface{}
		err   error
	)

	if "" == field {
		reply, err = r.Do(redisHLen, key)
	} else {
		reply, err = r.Do(redisHLen, key, field)
	}

	return redis.Int64(reply, err)
}

func (r *redigoImpl) Get(key string) (string, error) {
	reply, err := redis.String(r.Do(redisGet, key))
	if err == redis.ErrNil && r.UseCommonErr {
		return reply, ErrNil
	}
	return reply, err
}

func (r *redigoImpl) IncrBy(key string, incr int64) (int64, error) {
	reply, err := r.Do(redisIncrBy, key, incr)
	return redis.Int64(reply, err)
}

func (r *redigoImpl) Del(key ...string) error {
	if len(key) == 0 {
		return ErrInsufficientArgument
	}
	_, err := r.Do(redisDel, redis.Args{}.AddFlat(key)...)
	return err
}

func (r *redigoImpl) ZAdd(key, member string, score int) error {
	_, err := r.Do(redisZAdd, key, score, member)
	return err
}

func (r *redigoImpl) ZAddXX(key, member string, score int) error {
	_, err := r.Do(redisZAdd, key, redisXX, score, member)
	return err
}

func (r *redigoImpl) ZAddNX(key, member string, score int64) (int64, error) {
	return redis.Int64(r.Do(redisZAdd, key, redisNX, score, member))
}

func (r *redigoImpl) ZAddINCR(key, member string, score int) error {
	_, err := r.Do(redisZAdd, key, redisINCR, score, member)
	return err
}

func (r *redigoImpl) ZCard(key string) (int64, error) {
	return redis.Int64(r.Do(redisZCard, key))
}

func (r *redigoImpl) ZRange(key string, start, stop int) ([]string, error) {
	reply, err := r.Do(redisZRange, key, start, stop)
	return redis.Strings(reply, err)
}

func (r *redigoImpl) ZRevRange(key string, start, stop int) ([]string, error) {
	reply, err := r.Do(redisZRevRange, key, start, stop)
	return redis.Strings(reply, err)
}

func (r *redigoImpl) ZRangeByScore(key string, min, max, offset, count int) ([]string, error) {
	reply, err := r.Do(redisZRangeByScore, key, min, max, redisLimit, offset, count)
	return redis.Strings(reply, err)
}

func (r *redigoImpl) ZRevRangeByScore(key string, max, min, offset, count int) ([]string, error) {
	reply, err := r.Do(redisZRevRangeByScore, key, max, min, redisLimit, offset, count)
	return redis.Strings(reply, err)
}

func (r *redigoImpl) ZRank(key, member string) (int64, error) {
	reply, err := r.Do(redisZRank, key, member)
	return redis.Int64(reply, err)
}

func (r *redigoImpl) ZRevRank(key, member string) (int64, error) {
	reply, err := r.Do(redisZRevRank, key, member)
	return redis.Int64(reply, err)
}

func (r *redigoImpl) ZScore(key, member string) (int64, error) {
	reply, err := r.Do(redisZScore, key, member)
	return redis.Int64(reply, err)
}

func (r *redigoImpl) ZCount(key string, min, max int) (int64, error) {
	reply, err := r.Do(redisZCount, key, min, max)
	return redis.Int64(reply, err)
}

func (r *redigoImpl) ZRemRangeByScore(key string, start, stop int) (int64, error) {
	return redis.Int64(r.Do(redisZRemRangeByScore, key, start, stop))
}

func (r *redigoImpl) SAdd(key, member string) (int64, error) {
	return redis.Int64(r.Do(redisSAdd, key, member))
}

func (r *redigoImpl) SCard(key string) (int64, error) {
	return redis.Int64(r.Do(redisSCard, key))
}

func (r *redigoImpl) SDiff(keys ...string) ([]string, error) {
	return redis.Strings(r.Do(redisSDiff, convertArrayStringsToArrayInterfaces(keys)...))
}

func (r *redigoImpl) SDiffStore(keys ...string) (int64, error) {
	return redis.Int64(r.Do(redisSDiffStore, convertArrayStringsToArrayInterfaces(keys)...))
}

func (r *redigoImpl) SInter(keys ...string) ([]string, error) {
	return redis.Strings(r.Do(redisSInter, convertArrayStringsToArrayInterfaces(keys)...))
}

func (r *redigoImpl) SInterStore(keys ...string) (int64, error) {
	return redis.Int64(r.Do(redisSInterStore, convertArrayStringsToArrayInterfaces(keys)...))
}

func (r *redigoImpl) SIsMember(keys, member string) (int64, error) {
	return redis.Int64(r.Do(redisSIsMember, keys, member))
}

func (r *redigoImpl) SMembers(key string) ([]string, error) {
	return redis.Strings(r.Do(redisSMembers, key))
}

func (r *redigoImpl) SMove(value, source, destination string) (int64, error) {
	return redis.Int64(r.Do(redisSMove, source, destination, value))
}

func (r *redigoImpl) SPop(key string, count int) ([]string, error) {
	return redis.Strings(r.Do(redisSPop, key, count))
}

func (r *redigoImpl) SRandMember(key string, count int) ([]string, error) {
	return redis.Strings(r.Do(redisSRandMember, key, count))
}

func (r *redigoImpl) SRem(key string, member string) (int64, error) {
	return redis.Int64(r.Do(redisSRem, key, member))
}

func (r *redigoImpl) SUnion(keys ...string) ([]string, error) {
	return redis.Strings(r.Do(redisSUnion, convertArrayStringsToArrayInterfaces(keys)...))
}

func (r *redigoImpl) SUnionStore(keys ...string) (int64, error) {
	return redis.Int64(r.Do(redisSUnionStore, convertArrayStringsToArrayInterfaces(keys)...))
}

func convertArrayStringsToArrayInterfaces(input []string) []interface{} {
	output := make([]interface{}, 0)
	for _, v := range input {
		output = append(output, v)
	}
	return output
}

func (r *redigoImpl) ZRem(key string, members ...string) (int64, error) {
	return redis.Int64(r.Do(redisZREM, redis.Args{}.Add(key).AddFlat(members)...))
}

func (r *redigoImpl) ZAddXXIncrBy(key, member string, incrValue int64) (int64, error) {
	reply, err := redis.Int64(r.Do(redisZAdd, key, redisXX, redisINCR, incrValue, member))
	if err == redis.ErrNil {
		return 0, ErrXX
	}
	return reply, err
}

func (r *redigoImpl) HExists(key, field string) (bool, error) {
	return redis.Bool(r.Do(redisHExists, key, field))
}

func (r *redigoImpl) HIncrBy(key, field string, incrValue int64) (int64, error) {
	return redis.Int64(r.Do(redisHIncrBy, key, field, incrValue))
}

func (r *redigoImpl) Expire(key string, ttl time.Duration) (int64, error) {
	return redis.Int64(r.Do(redisExpire, key, int64(ttl.Seconds())))
}

func (r *redigoImpl) TTL(key string) (int64, error) {
	reply, err := redis.Int64(r.Do(redisTTL, key))

	if reply == -2 {
		return reply, ErrNil
	}

	return reply, err
}

func (r *redigoImpl) Exists(key string) (bool, error) {
	return redis.Bool(r.Do(redisExists, key))
}

func (r *redigoImpl) ZRevRangeWithScore(key string, start, stop int64) (interface{}, error) {
	reply, err := r.Do(redisZRevRange, key, start, stop, redisWithScores)
	return redis.Strings(reply, err)
}

type GeoRadiusUnit string
type GeoRadiusSort string

const (
	GeoRadiusMeter     GeoRadiusUnit = "m"
	GeoRadiusKiloMeter GeoRadiusUnit = "km"
	GeoRadiusFeet      GeoRadiusUnit = "ft"
	GeoRadiusMile      GeoRadiusUnit = "mi"
	GeoRadiusAsc       GeoRadiusSort = "asc"
	GeoRadiusDesc      GeoRadiusSort = "desc"
)

// GeoPoint is geospatial information in point geometry object
type GeoPoint struct {
	Member    string
	Latitude  float64
	Longitude float64
}

// GeoLoc is geospatial information
type GeoLoc struct {
	Name                          string
	Longitude, Latitude, Distance float64
	GeoHash                       int64
}

// GeoRadiusQuery is used with GeoRadius
type GeoRadiusQuery struct {
	Radius float64
	// Can be m, km, ft, or mi. Default is km.
	Unit        GeoRadiusUnit
	WithCoord   bool
	WithDist    bool
	WithGeoHash bool
	Count       int
	// Can be ASC or DESC. Default is no sort order.
	Sort GeoRadiusSort
}

// GeoAdd Adds the specified geospatial items to the specified key
func (r *redigoImpl) GeoAdd(key string, geos ...*GeoPoint) (int64, error) {
	arg := []interface{}{
		key,
	}
	for _, g := range geos {
		arg = append(arg, g.Longitude, g.Latitude, g.Member)
	}

	reply, err := r.Do(redisGeoAdd, arg...)
	return redis.Int64(reply, err)
}

// GeoHash return valid Geohash strings representing the position of one or more elements in a sorted set value representing a geospatial index
func (r *redigoImpl) GeoHash(key string, members ...string) ([]string, error) {
	return redis.Strings(r.Do(redisGeoHash, redis.Args{}.Add(key).AddFlat(members)...))
}

// GeoRadius get items from given key within given radius from particular given point
func (r *redigoImpl) GeoRadius(key string, long, lat float64, q *GeoRadiusQuery) ([]*GeoLoc, error) {
	args := []interface{}{
		key,
		long,
		lat,
	}

	args = append(args, q.Radius)
	switch q.Unit {
	case GeoRadiusMeter, GeoRadiusKiloMeter, GeoRadiusFeet, GeoRadiusMile:
		args = append(args, string(q.Unit))
	default:
		args = append(args, string(GeoRadiusKiloMeter))
	}
	if q.WithCoord {
		args = append(args, "WITHCOORD")
	}
	if q.WithDist {
		args = append(args, "WITHDIST")
	}
	if q.WithGeoHash {
		args = append(args, "WITHHASH")
	}
	if q.Count > 0 {
		args = append(args, "count", q.Count)
	}
	switch q.Sort {
	case GeoRadiusAsc, GeoRadiusDesc:
		args = append(args, string(q.Sort))
	}

	reply, err := r.Do(redisGeoRadius, args...)
	locs, err := redis.Values(reply, err)
	result := []*GeoLoc{}

	for _, loc := range locs {
		switch v := loc.(type) {
		case []byte:
			result = append(result, &GeoLoc{
				Name: string(v),
			})
		case []interface{}:
			locAttr := v
			var (
				name                string
				distance, lat, long float64
				geohash             int64
			)

			if len(locAttr) > 0 {
				name = fmt.Sprintf("%s", locAttr[0])
			}

			for i := 1; i < len(locAttr); i++ {
				a := locAttr[i]

				switch vv := a.(type) {
				case []byte: // should be distance b/c name already taken
					distance, _ = strconv.ParseFloat(string(vv), 64)
				case int64: // should be geohash
					geohash = vv
				case []interface{}: // should be long, lat
					if len(vv) < 2 {
						return nil, fmt.Errorf("unexpected return value")
					}

					longBytes, ok := vv[0].([]byte)
					if !ok {
						return nil, fmt.Errorf("unexpected return value")
					}

					latBytes, ok := vv[1].([]byte)
					if !ok {
						return nil, fmt.Errorf("unexpected return value")
					}

					long, _ = strconv.ParseFloat(string(longBytes), 64)
					lat, _ = strconv.ParseFloat(string(latBytes), 64)
				default:
					return nil, fmt.Errorf("got %T, expected []byte or int64 or []interface{}", v)
				}

			}

			result = append(result, &GeoLoc{
				Name:      name,
				Distance:  distance,
				GeoHash:   geohash,
				Longitude: long,
				Latitude:  lat,
			})

		default:
			return nil, fmt.Errorf("got %T, expected string or []interface{}", v)
		}
	}

	return result, err
}

func (r *redigoImpl) MGet(keys []string) ([]string, error) {
	c := r.GetConn()
	defer c.Close()

	iKeys := make([]interface{}, len(keys))
	for k := range keys {
		iKeys[k] = keys[k]
	}

	res, err := redis.Values(c.Do(redisMGet, iKeys...))
	if err != nil {
		return nil, err
	}

	rsp := make([]string, len(keys))
	for k := range keys {
		if res[k] != nil {
			rsp[k] = string(res[k].([]byte))
		}
	}

	return rsp, err
}

func (r *redigoImpl) MSet(values map[string]string) error {
	c := r.GetConn()
	defer c.Close()

	val := make([]interface{}, 0)
	for key := range values {
		val = append(val, key, values[key])
	}

	_, err := c.Do(redisMSet, val...)
	return err
}

func (r *redigoImpl) MSetEx(values map[string]string, ttl time.Duration) error {
	c := r.GetConn()
	defer c.Close()

	c.Send(redisMulti)
	for key := range values {
		c.Send(redisSet, key, values[key], redisEx, int(ttl.Seconds()))
	}
	_, err := c.Do(redisExec)
	return err
}

func (r *redigoImpl) LLen(key string) (int64, error) {
	c := r.GetConn()
	defer c.Close()

	val, err := c.Do(redisLLen, key)
	if err != nil {
		return 0, err
	}

	length, _ := val.(int64)
	return length, nil
}

func (r *redigoImpl) LPop(key string, count int) ([]string, error) {
	c := r.GetConn()
	defer c.Close()

	res, err := redis.Values(c.Do(redisLPop, key, count))
	if err != nil {
		return nil, err
	}

	rsp := make([]string, 0, len(res))
	for i := 0; i < len(res); i++ {
		rsp = append(rsp, string(res[i].([]byte)))
	}

	return rsp, err
}

func (r *redigoImpl) LPush(key string, values []string) (int64, error) {
	c := r.GetConn()
	defer c.Close()

	res, err := c.Do(redisLPush, redis.Args{}.Add(key).AddFlat(values)...)
	if err != nil {
		return 0, err
	}

	length, _ := res.(int64)
	return length, nil
}

func (r *redigoImpl) LPushX(key string, values []string) (int64, error) {
	c := r.GetConn()
	defer c.Close()

	res, err := c.Do(redisLPushX, redis.Args{}.Add(key).AddFlat(values)...)
	if err != nil {
		return 0, err
	}

	length, _ := res.(int64)
	if length == 0 {
		return length, ErrNil
	}
	return length, nil
}

func (r *redigoImpl) RPop(key string, count int) ([]string, error) {
	c := r.GetConn()
	defer c.Close()

	res, err := redis.Values(c.Do(redisRPop, key, count))
	if err != nil {
		return nil, err
	}

	rsp := make([]string, 0, len(res))
	for i := 0; i < len(res); i++ {
		rsp = append(rsp, string(res[i].([]byte)))
	}

	return rsp, err
}

func (r *redigoImpl) RPush(key string, values []string) (int64, error) {
	c := r.GetConn()
	defer c.Close()

	res, err := c.Do(redisRPush, redis.Args{}.Add(key).AddFlat(values)...)
	if err != nil {
		return 0, err
	}

	length, _ := res.(int64)
	return length, nil
}

func (r *redigoImpl) RPushX(key string, values []string) (int64, error) {
	c := r.GetConn()
	defer c.Close()

	res, err := c.Do(redisRPushX, redis.Args{}.Add(key).AddFlat(values)...)
	if err != nil {
		return 0, err
	}

	length, _ := res.(int64)
	if length == 0 {
		return length, ErrNil
	}
	return length, nil
}
