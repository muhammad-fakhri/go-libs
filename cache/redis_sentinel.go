package cache

import (
	"errors"
	"strconv"
	"time"

	goredis "github.com/go-redis/redis/v7"
	redigo "github.com/gomodule/redigo/redis"
)

type redisSentinelImpl struct {
	client *goredis.Client
}

func (r *redisSentinelImpl) GetConn() Conn {
	return r
}

func (r *redisSentinelImpl) Close() error {
	return r.client.Close()
}

func (r *redisSentinelImpl) Err() error {
	return ErrNotSupported
}

func (r *redisSentinelImpl) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return r.client.Do(commandName, args).Result()
}

func (r *redisSentinelImpl) Send(commandName string, args ...interface{}) error {
	return ErrNotSupported
}

func (r *redisSentinelImpl) Flush() error {
	return ErrNotSupported
}

func (r *redisSentinelImpl) Receive() (reply interface{}, err error) {
	return nil, ErrNotSupported
}

var (
	ErrFailInitialize = errors.New("failed to initiate conn to redis sentinel")
)

func newRedisSentinel(opt *FailoverOptions) (*redisSentinelImpl, error) {
	sentinel := redisSentinelImpl{}

	sentinel.client = goredis.NewFailoverClient(&goredis.FailoverOptions{
		MasterName:         opt.MasterName,
		SentinelAddrs:      opt.SentinelAddrs,
		OnConnect:          opt.OnConnect,
		Password:           opt.Password,
		DB:                 opt.DB,
		MaxRetries:         opt.MaxRetries,
		MinRetryBackoff:    opt.MinRetryBackoff,
		MaxRetryBackoff:    opt.MaxRetryBackoff,
		DialTimeout:        opt.DialTimeout,
		ReadTimeout:        opt.ReadTimeout,
		WriteTimeout:       opt.WriteTimeout,
		PoolSize:           opt.PoolSize,
		MinIdleConns:       opt.MinIdleConns,
		MaxConnAge:         opt.MaxConnAge,
		PoolTimeout:        opt.PoolTimeout,
		IdleTimeout:        opt.IdleTimeout,
		IdleCheckFrequency: opt.IdleCheckFrequency,
		TLSConfig:          opt.TLSConfig,
	})

	if sentinel.client == nil {
		return nil, ErrFailInitialize
	}

	_, err := sentinel.client.Ping().Result()

	return &sentinel, err
}

func (r *redisSentinelImpl) Set(key, value string, ttl time.Duration) error {
	var err error

	if 0 >= ttl {
		_, err = r.client.Set(key, value, 0).Result()
	} else {
		_, err = r.client.Set(key, value, ttl).Result()
	}
	return err
}

func (r *redisSentinelImpl) Get(key string) (string, error) {

	result, err := r.client.Get(key).Result()

	if err != nil && err == goredis.Nil {
		return result, ErrNil
	}
	return result, err
}

func (r *redisSentinelImpl) Del(key ...string) error {
	if len(key) == 0 {
		return ErrInsufficientArgument
	}
	_, err := r.client.Del(key...).Result()
	return err
}

func (r *redisSentinelImpl) ErrorOnCacheMiss() error {
	return ErrNil
}

func (r *redisSentinelImpl) HSet(key, field, value string, ttl time.Duration) error {
	var err error

	_, err = r.client.HSet(key, field, value).Result()

	if err != nil {
		return err
	}

	if ttl > 0 {
		_, err = r.client.Expire(key, ttl).Result()
	}

	return err
}

func (r *redisSentinelImpl) HSetNX(key, field, value string, ttl time.Duration) error {
	var err error

	reply, err := r.client.HSetNX(key, field, value).Result()
	if err != nil {
		return err
	}
	if reply == false {
		return ErrHNX
	}

	if ttl > 0 {
		_, err = r.client.Expire(key, ttl).Result()
	}

	return err
}

func (r *redisSentinelImpl) HMSet(key string, fieldsMap map[string]string, ttl time.Duration) error {

	values := map[string]interface{}{}

	for index, value := range fieldsMap {
		values[index] = value
	}

	_, err := r.client.HMSet(key, values).Result()

	if ttl > 0 {
		_, err = r.client.Expire(key, ttl).Result()
	}

	return err
}

func (r *redisSentinelImpl) HGet(key, field string) (string, error) {
	result, err := r.client.HGet(key, field).Result()
	if err != nil && err == goredis.Nil {
		return result, ErrNil
	}
	return result, err
}

func (r *redisSentinelImpl) HMGet(key string, fields ...string) ([]string, error) {
	var keys []string
	out, err := r.client.HMGet(key, fields...).Result()

	if err != nil {
		return keys, err
	}

	for _, result := range out {
		keys = append(keys, result.(string))
	}
	return keys, err
}

func (r *redisSentinelImpl) HDel(key string, fields ...string) (int64, error) {
	return r.client.HDel(key, fields...).Result()
}

func (r *redisSentinelImpl) HKeys(key string) ([]string, error) {
	return r.client.HKeys(key).Result()
}

func (r *redisSentinelImpl) HVals(key string) ([]string, error) {
	return r.client.HVals(key).Result()
}

func (r *redisSentinelImpl) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(key).Result()
}

func (r *redisSentinelImpl) HExists(key, field string) (bool, error) {
	return r.client.HExists(key, field).Result()
}

func (r *redisSentinelImpl) HIncrBy(key, field string, incrValue int64) (int64, error) {
	return r.client.HIncrBy(key, field, incrValue).Result()
}

func (r *redisSentinelImpl) ErrorOnHashCacheMiss() error {
	return ErrNil
}

func (r *redisSentinelImpl) SetNX(key, value string, ttl time.Duration) error {
	var (
		err   error
		reply interface{}
	)

	reply, err = r.client.SetNX(key, value, ttl).Result()

	if nil != err {
		return err
	}
	if false == reply {
		return ErrNX
	}
	return err
}

func (r *redisSentinelImpl) ScanKeys(pattern string) ([]string, error) {
	iter := 0
	keys := []string{}
	for {
		arr, err := redigo.Values(r.client.Do(redisScan, iter, redisMatch, pattern).Result())
		if err != nil {
			return nil, err
		}

		iter, _ = redigo.Int(arr[0], nil)
		k, _ := redigo.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}
	return keys, nil
}

func (r *redisSentinelImpl) IncrBy(key string, incr int64) (int64, error) {
	return r.client.IncrBy(key, incr).Result()
}

func (r *redisSentinelImpl) ZAdd(key, member string, score int) error {
	_, err := r.client.ZAdd(key, &goredis.Z{Member: member, Score: float64(score)}).Result()
	return err
}

func (r *redisSentinelImpl) ZAddXX(key, member string, score int) error {
	_, err := r.client.ZAddXX(key, &goredis.Z{Member: member, Score: float64(score)}).Result()
	return err
}

func (r *redisSentinelImpl) ZAddNX(key, member string, score int64) (int64, error) {
	return r.client.ZAddNX(key, &goredis.Z{Member: member, Score: float64(score)}).Result()
}

func (r *redisSentinelImpl) ZAddINCR(key, member string, score int) error {
	_, err := r.client.ZIncr(key, &goredis.Z{Member: member, Score: float64(score)}).Result()
	return err
}

func (r *redisSentinelImpl) ZCard(key string) (int64, error) {
	return r.client.ZCard(key).Result()
}

func (r *redisSentinelImpl) ZRange(key string, start, stop int) ([]string, error) {
	return r.client.ZRange(key, int64(start), int64(stop)).Result()
}

func (r *redisSentinelImpl) ZRevRange(key string, start, stop int) ([]string, error) {
	return r.client.ZRevRange(key, int64(start), int64(stop)).Result()
}

func (r *redisSentinelImpl) ZRangeByScore(key string, min, max, offset, count int) ([]string, error) {
	return r.client.ZRangeByScore(key, &goredis.ZRangeBy{Min: strconv.Itoa(min), Max: strconv.Itoa(max),
		Offset: int64(offset), Count: int64(count)}).Result()
}

func (r *redisSentinelImpl) ZRevRangeByScore(key string, max, min, offset, count int) ([]string, error) {
	return r.client.ZRevRangeByScore(key, &goredis.ZRangeBy{Min: strconv.Itoa(min),
		Max: strconv.Itoa(max), Offset: int64(offset), Count: int64(count)}).Result()
}

func (r *redisSentinelImpl) ZRevRangeWithScore(key string, start, stop int64) (interface{}, error) {
	return r.client.ZRevRangeWithScores(key, start, stop).Result()
}

func (r *redisSentinelImpl) ZRank(key, member string) (int64, error) {
	return r.client.ZRank(key, member).Result()
}

func (r *redisSentinelImpl) ZRevRank(key, member string) (int64, error) {
	return r.client.ZRevRank(key, member).Result()
}

func (r *redisSentinelImpl) ZScore(key, member string) (int64, error) {
	result, err := r.client.ZScore(key, member).Result()
	return int64(result), err
}

func (r *redisSentinelImpl) ZCount(key string, min, max int) (int64, error) {
	return r.client.ZCount(key, strconv.Itoa(min), strconv.Itoa(max)).Result()
}

func (r *redisSentinelImpl) ZRemRangeByScore(key string, start, stop int) (int64, error) {
	return r.client.ZRemRangeByScore(key, strconv.Itoa(start), strconv.Itoa(stop)).Result()
}

func (r *redisSentinelImpl) SAdd(key, member string) (int64, error) {
	return r.client.SAdd(key, member).Result()
}

func (r *redisSentinelImpl) SCard(key string) (int64, error) {
	return r.client.SCard(key).Result()
}

func (r *redisSentinelImpl) SDiff(keys ...string) ([]string, error) {
	return r.client.SDiff(keys...).Result()
}

func (r *redisSentinelImpl) SDiffStore(keys ...string) (int64, error) {
	if len(keys) < 2 {
		return 0, errors.New("ERR wrong number of arguments for 'sdiffstore' command") // same as redis/redigo error
	}

	return r.client.SDiffStore(keys[0], keys[1:]...).Result()
}

func (r *redisSentinelImpl) SInter(keys ...string) ([]string, error) {
	return r.client.SInter(keys...).Result()
}

func (r *redisSentinelImpl) SInterStore(keys ...string) (int64, error) {
	if len(keys) < 2 {
		return 0, errors.New("ERR wrong number of arguments for 'sdiffstore' command") // same as redis/redigo error
	}

	return r.client.SInterStore(keys[0], keys[1:]...).Result()
}

func (r *redisSentinelImpl) SIsMember(keys, member string) (int64, error) {
	result, err := r.client.SIsMember(keys, member).Result()

	if result {
		return 1, err
	}

	return 0, err
}

func (r *redisSentinelImpl) SMembers(key string) ([]string, error) {
	return r.client.SMembers(key).Result()
}

func (r *redisSentinelImpl) SMove(value, source, destination string) (int64, error) {
	res, err := r.client.SMove(source, destination, value).Result()

	result := int64(0)

	if res {
		result = 1
	}

	return result, err
}

func (r *redisSentinelImpl) SPop(key string, count int) ([]string, error) {
	return r.client.SPopN(key, int64(count)).Result()
}

func (r *redisSentinelImpl) SRandMember(key string, count int) ([]string, error) {
	return r.client.SRandMemberN(key, int64(count)).Result()
}

func (r *redisSentinelImpl) SRem(key string, member string) (int64, error) {
	return r.client.SRem(key, member).Result()
}

func (r *redisSentinelImpl) SUnion(keys ...string) ([]string, error) {
	return r.client.SUnion(keys...).Result()
}

func (r *redisSentinelImpl) SUnionStore(keys ...string) (int64, error) {
	args := append([]string{redisSUnionStore}, keys...)

	return r.client.Do(convertArrayStringsToArrayInterfaces(args)...).Int64()
}

func (r *redisSentinelImpl) ZRem(key string, members ...string) (int64, error) {
	return r.client.ZRem(key, members).Result()
}

func (r *redisSentinelImpl) ZAddXXIncrBy(key, member string, incrValue int64) (int64, error) {
	result, err := r.client.ZIncrBy(key, float64(incrValue), member).Result()
	return int64(result), err
}

func (r *redisSentinelImpl) Expire(key string, ttl time.Duration) (int64, error) {
	result, err := r.client.Expire(key, ttl).Result()

	if !result {
		return 0, err
	}
	return 1, nil
}

func (r *redisSentinelImpl) TTL(key string) (int64, error) {
	duration, err := r.client.TTL(key).Result()
	ttl := int64(duration.Seconds())

	if ttl == -2 { // key not found
		return ttl, ErrNil
	}

	return ttl, err
}

func (r *redisSentinelImpl) Exists(key string) (bool, error) {
	result, err := r.client.Exists(key).Result()

	if result == 1 {
		return true, err
	}

	return false, err
}

func (r *redisSentinelImpl) IncrXX(key string, value int64) (reply int64, err error) {
	if exists, _ := r.Exists(key); !exists {
		return 0, ErrXX
	}

	return r.client.IncrBy(key, value).Result()
}

func (r *redisSentinelImpl) DecrWithLimit(key string, value, lowerBound int64) (reply int64, err error) {
	decrWithLimitScript := goredis.NewScript(`
		local key = KEYS[1]
		local decrement = tonumber(ARGV[1])
		local lb = tonumber(ARGV[2])
		local cnt = redis.call('get', key) or 0
		cnt = cnt - decrement
		if (cnt >= lb ) then
			redis.call('decrby', key, decrement)
			return cnt
		else 
			return lb - 1
		end
	`)

	args := []string{strconv.Itoa(int(value)), strconv.Itoa(int(lowerBound))}

	reply, err = decrWithLimitScript.Run(r.client, []string{key}, convertArrayStringsToArrayInterfaces(args)...).Int64()
	if err != nil {
		return lowerBound - 1, err
	}

	if reply < lowerBound {
		return lowerBound - 1, ErrLimitExceeded
	}

	return
}

func (r *redisSentinelImpl) HGetSet(key, field, value, prevValue string, ttl time.Duration) error {
	IncrByXX := goredis.NewScript(`
		local key = KEYS[1]
		local column = ARGV[1]
		local expireSecond = ARGV[2]
		local prevValue = redis.call('HGET', key, column)
		
		if (prevValue==ARGV[4]) then 
			redis.call('HSET', key, column,ARGV[3])
			redis.call('expire', key, expireSecond)
			return ARGV[3]
		else
			return 'errValueInvalid'
		end
	`)

	result, err := IncrByXX.Run(r.client, []string{key}, field, int(ttl), value, prevValue).Result()

	if result == valueInvalid {
		return ErrValueInvalid
	}

	return err
}

func (thisCluster *redisSentinelImpl) ZAddToFixed(key, member string, score, maxSize int) (reply int64, err error) {
	return 0, ErrNotSupported
}

func (thisCluster *redisSentinelImpl) GeoAdd(key string, geos ...*GeoPoint) (int64, error) {
	goredisGeoLoc := []*goredis.GeoLocation{}
	for _, g := range geos {
		goredisGeoLoc = append(goredisGeoLoc, &goredis.GeoLocation{
			Name:      g.Member,
			Longitude: g.Longitude,
			Latitude:  g.Latitude,
		})
	}

	return thisCluster.client.GeoAdd(key, goredisGeoLoc...).Result()
}

func (thisCluster *redisSentinelImpl) GeoHash(key string, members ...string) ([]string, error) {
	return thisCluster.client.GeoHash(key, members...).Result()
}

func (thisCluster *redisSentinelImpl) GeoRadius(key string, long, lat float64, q *GeoRadiusQuery) ([]*GeoLoc, error) {
	goredisGeoRadiusQuery := &goredis.GeoRadiusQuery{
		Radius:      q.Radius,
		Unit:        string(q.Unit),
		WithCoord:   q.WithCoord,
		WithDist:    q.WithDist,
		WithGeoHash: q.WithGeoHash,
		Count:       q.Count,
		Sort:        string(q.Sort),
	}
	result := []*GeoLoc{}
	goredisGeoLoc, err := thisCluster.client.GeoRadius(key, long, lat, goredisGeoRadiusQuery).Result()
	if err != nil {
		return result, err
	}

	for _, g := range goredisGeoLoc {
		result = append(result, &GeoLoc{
			Name:      g.Name,
			Distance:  g.Dist,
			GeoHash:   g.GeoHash,
			Longitude: g.Longitude,
			Latitude:  g.Latitude,
		})
	}

	return result, nil
}

func (thisCluster *redisSentinelImpl) MGet(keys []string) ([]string, error) {
	return nil, ErrNotSupported
}

func (thisCluster *redisSentinelImpl) MSet(values map[string]string) error {
	return ErrNotSupported
}

func (thisCluster *redisSentinelImpl) MSetEx(values map[string]string, ttl time.Duration) error {
	return ErrNotSupported
}

func (thisCluster *redisSentinelImpl) LLen(key string) (int64, error) {
	return 0, ErrNotSupported
}

func (thisCluster *redisSentinelImpl) LPop(key string, count int) ([]string, error) {
	return nil, ErrNotSupported
}

func (thisCluster *redisSentinelImpl) LPush(key string, values []string) (int64, error) {
	return 0, ErrNotSupported
}

func (thisCluster *redisSentinelImpl) LPushX(key string, values []string) (int64, error) {
	return 0, ErrNotSupported
}

func (thisCluster *redisSentinelImpl) RPop(key string, count int) ([]string, error) {
	return nil, ErrNotSupported
}

func (thisCluster *redisSentinelImpl) RPush(key string, values []string) (int64, error) {
	return 0, ErrNotSupported
}

func (thisCluster *redisSentinelImpl) RPushX(key string, values []string) (int64, error) {
	return 0, ErrNotSupported
}
