package cache

import (
	"errors"
	"strconv"
	"time"

	gocluster "github.com/go-redis/redis/v7"
)

var (
	ErrNotSupported = errors.New("command not supported on cluster")
)

type goRedisClusterImpl struct {
	currentClusterClient *gocluster.ClusterClient
}

func (thisCluster *goRedisClusterImpl) GetConn() Conn {
	return thisCluster
}

func newRedisCluster(currentNodes *ConfigCacheCluster) (*goRedisClusterImpl, error) {
	currentCluster := goRedisClusterImpl{}
	currentCluster.currentClusterClient = gocluster.NewClusterClient(&gocluster.ClusterOptions{
		Addrs:          currentNodes.Addrs,
		MaxRedirects:   currentNodes.MaxRedirects,
		ReadOnly:       currentNodes.ReadOnly,
		RouteByLatency: currentNodes.RouteByLatency,
		RouteRandomly:  currentNodes.RouteRandomly,
		ReadTimeout:    currentNodes.ReadTimeout,
		WriteTimeout:   currentNodes.WriteTimeout,
		PoolSize:       currentNodes.PoolSize,
		DialTimeout:    currentNodes.DialTimeout,
		IdleTimeout:    currentNodes.IdleTimeout,
	})

	_, err := currentCluster.currentClusterClient.Ping().Result()

	return &currentCluster, err
}
func (thisCluster *goRedisClusterImpl) Set(key, value string, ttl time.Duration) error {
	var err error

	if 0 >= ttl {
		_, err = thisCluster.currentClusterClient.Set(key, value, 0).Result()
	} else {
		_, err = thisCluster.currentClusterClient.Set(key, value, ttl).Result()
	}
	return err
}

func (thisCluster *goRedisClusterImpl) Get(key string) (string, error) {
	result, err := thisCluster.currentClusterClient.Get(key).Result()
	if err == gocluster.Nil {
		return result, ErrNil
	}
	return result, err
}

func (thisCluster *goRedisClusterImpl) Del(key ...string) error {
	if len(key) == 0 {
		return ErrInsufficientArgument
	}
	_, err := thisCluster.currentClusterClient.Del(key...).Result()
	return err
}

func (thisCluster *goRedisClusterImpl) ErrorOnCacheMiss() error {
	return ErrNil
}

func (thisCluster *goRedisClusterImpl) HSet(key, field, value string, ttl time.Duration) error {
	var err error

	_, err = thisCluster.currentClusterClient.HSet(key, field, value).Result()

	if err != nil {
		return err
	}

	if ttl > 0 {
		_, err = thisCluster.currentClusterClient.Expire(key, ttl).Result()
	}

	return err
}

func (thisCluster *goRedisClusterImpl) HSetNX(key, field, value string, ttl time.Duration) error {
	var err error

	reply, err := thisCluster.currentClusterClient.HSetNX(key, field, value).Result()
	if err != nil {
		return err
	}
	if reply == false {
		return ErrHNX
	}

	if ttl > 0 {
		_, err = thisCluster.currentClusterClient.Expire(key, ttl).Result()
	}

	return err
}

func (thisCluster *goRedisClusterImpl) HMSet(key string, fieldsMap map[string]string, ttl time.Duration) error {

	values := map[string]interface{}{}

	for index, value := range fieldsMap {
		values[index] = value
	}

	_, err := thisCluster.currentClusterClient.HMSet(key, values).Result()

	if ttl > 0 {
		_, err = thisCluster.currentClusterClient.Expire(key, ttl).Result()
	}

	return err
}

func (thisCluster *goRedisClusterImpl) HGet(key, field string) (string, error) {
	result, err := thisCluster.currentClusterClient.HGet(key, field).Result()
	if err == gocluster.Nil {
		return result, ErrNil
	}
	return result, err
}

func (thisCluster *goRedisClusterImpl) HMGet(key string, fields ...string) ([]string, error) {
	var keys []string
	out, err := thisCluster.currentClusterClient.HMGet(key, fields...).Result()

	if err != nil {
		return keys, err
	}

	for _, result := range out {
		keys = append(keys, result.(string))
	}
	return keys, err
}

func (thisCluster *goRedisClusterImpl) HDel(key string, fields ...string) (int64, error) {
	return thisCluster.currentClusterClient.HDel(key, fields...).Result()
}

func (thisCluster *goRedisClusterImpl) HKeys(key string) ([]string, error) {
	return thisCluster.currentClusterClient.HKeys(key).Result()
}

func (thisCluster *goRedisClusterImpl) HVals(key string) ([]string, error) {
	return thisCluster.currentClusterClient.HVals(key).Result()
}

func (thisCluster *goRedisClusterImpl) HGetAll(key string) (map[string]string, error) {
	return thisCluster.currentClusterClient.HGetAll(key).Result()
}

func (thisCluster *goRedisClusterImpl) HExists(key, field string) (bool, error) {
	return thisCluster.currentClusterClient.HExists(key, field).Result()
}

func (thisCluster *goRedisClusterImpl) HIncrBy(key, field string, incrValue int64) (int64, error) {
	return thisCluster.currentClusterClient.HIncrBy(key, field, incrValue).Result()
}

func (thisCluster *goRedisClusterImpl) ErrorOnHashCacheMiss() error {
	return ErrNil
}

func (thisCluster *goRedisClusterImpl) SetNX(key, value string, ttl time.Duration) error {
	var (
		err   error
		reply interface{}
	)

	reply, err = thisCluster.currentClusterClient.SetNX(key, value, ttl).Result()

	if nil != err {
		return err
	}
	if false == reply {
		return ErrNX
	}
	return err
}

func (thisCluster *goRedisClusterImpl) ScanKeys(pattern string) ([]string, error) {
	return nil, ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) IncrBy(key string, incr int64) (int64, error) {
	return thisCluster.currentClusterClient.IncrBy(key, incr).Result()
}

func (thisCluster *goRedisClusterImpl) ZAdd(key, member string, score int) error {
	_, err := thisCluster.currentClusterClient.ZAdd(key, &gocluster.Z{Member: member, Score: float64(score)}).Result()
	return err
}

func (thisCluster *goRedisClusterImpl) ZAddXX(key, member string, score int) error {
	_, err := thisCluster.currentClusterClient.ZAddXX(key, &gocluster.Z{Member: member, Score: float64(score)}).Result()
	return err
}

func (thisCluster *goRedisClusterImpl) ZAddNX(key, member string, score int64) (int64, error) {
	return thisCluster.currentClusterClient.ZAddNX(key, &gocluster.Z{Member: member, Score: float64(score)}).Result()
}

func (thisCluster *goRedisClusterImpl) ZAddINCR(key, member string, score int) error {
	_, err := thisCluster.currentClusterClient.ZIncr(key, &gocluster.Z{Member: member, Score: float64(score)}).Result()
	return err
}

func (thisCluster *goRedisClusterImpl) ZCard(key string) (int64, error) {
	return thisCluster.currentClusterClient.ZCard(key).Result()
}

func (thisCluster *goRedisClusterImpl) ZRange(key string, start, stop int) ([]string, error) {
	return thisCluster.currentClusterClient.ZRange(key, int64(start), int64(stop)).Result()
}

func (thisCluster *goRedisClusterImpl) ZRevRange(key string, start, stop int) ([]string, error) {
	return thisCluster.currentClusterClient.ZRevRange(key, int64(start), int64(stop)).Result()
}

func (thisCluster *goRedisClusterImpl) ZRangeByScore(key string, min, max, offset, count int) ([]string, error) {
	return thisCluster.currentClusterClient.ZRangeByScore(key, &gocluster.ZRangeBy{Min: strconv.Itoa(min), Max: strconv.Itoa(max),
		Offset: int64(offset), Count: int64(count)}).Result()
}

func (thisCluster *goRedisClusterImpl) ZRevRangeByScore(key string, max, min, offset, count int) ([]string, error) {
	return thisCluster.currentClusterClient.ZRevRangeByScore(key, &gocluster.ZRangeBy{Min: strconv.Itoa(min),
		Max: strconv.Itoa(max), Offset: int64(offset), Count: int64(count)}).Result()
}

func (thisCluster *goRedisClusterImpl) ZRevRangeWithScore(key string, start, stop int64) (interface{}, error) {
	return thisCluster.currentClusterClient.ZRevRangeWithScores(key, start, stop).Result()
}

func (thisCluster *goRedisClusterImpl) ZRank(key, member string) (int64, error) {
	return thisCluster.currentClusterClient.ZRank(key, member).Result()
}

func (thisCluster *goRedisClusterImpl) ZRevRank(key, member string) (int64, error) {
	return thisCluster.currentClusterClient.ZRevRank(key, member).Result()
}

func (thisCluster *goRedisClusterImpl) ZScore(key, member string) (int64, error) {
	result, err := thisCluster.currentClusterClient.ZScore(key, member).Result()
	return int64(result), err
}

func (thisCluster *goRedisClusterImpl) ZCount(key string, min, max int) (int64, error) {
	return thisCluster.currentClusterClient.ZCount(key, strconv.Itoa(min), strconv.Itoa(max)).Result()
}

func (thisCluster *goRedisClusterImpl) ZRemRangeByScore(key string, start, stop int) (int64, error) {
	return thisCluster.currentClusterClient.ZRemRangeByScore(key, strconv.Itoa(start), strconv.Itoa(stop)).Result()
}

func (thisCluster *goRedisClusterImpl) SAdd(key, member string) (int64, error) {
	return thisCluster.currentClusterClient.SAdd(key, member).Result()
}

func (thisCluster *goRedisClusterImpl) SCard(key string) (int64, error) {
	return thisCluster.currentClusterClient.SCard(key).Result()
}

func (thisCluster *goRedisClusterImpl) SDiff(keys ...string) ([]string, error) {
	return nil, ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) SDiffStore(keys ...string) (int64, error) {
	return 0, ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) SInter(keys ...string) ([]string, error) {
	return nil, ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) SInterStore(keys ...string) (int64, error) {
	return 0, ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) SIsMember(keys, member string) (int64, error) {
	result, err := thisCluster.currentClusterClient.SIsMember(keys, member).Result()
	if result {
		return 1, err
	}
	return 0, err
}

func (thisCluster *goRedisClusterImpl) SMembers(key string) ([]string, error) {
	return thisCluster.currentClusterClient.SMembers(key).Result()
}

func (thisCluster *goRedisClusterImpl) SMove(value, source, destination string) (int64, error) {
	return 0, ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) SPop(key string, count int) ([]string, error) {
	return thisCluster.currentClusterClient.SPopN(key, int64(count)).Result()
}

func (thisCluster *goRedisClusterImpl) SRandMember(key string, count int) ([]string, error) {
	return thisCluster.currentClusterClient.SRandMemberN(key, int64(count)).Result()
}

func (thisCluster *goRedisClusterImpl) SRem(key string, member string) (int64, error) {
	return thisCluster.currentClusterClient.SRem(key, member).Result()
}

func (thisCluster *goRedisClusterImpl) SUnion(keys ...string) ([]string, error) {
	return nil, ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) SUnionStore(keys ...string) (int64, error) {
	return 0, ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) ZRem(key string, members ...string) (int64, error) {
	return thisCluster.currentClusterClient.ZRem(key, members).Result()
}

func (thisCluster *goRedisClusterImpl) ZAddXXIncrBy(key, member string, incrValue int64) (int64, error) {
	result, err := thisCluster.currentClusterClient.ZIncrBy(key, float64(incrValue), member).Result()
	return int64(result), err
}

func (thisCluster *goRedisClusterImpl) TTL(key string) (int64, error) {
	duration, err := thisCluster.currentClusterClient.TTL(key).Result()
	ttl := int64(duration.Seconds())

	if ttl == -2 { // key not found
		return ttl, ErrNil
	}

	return ttl, err
}

func (thisCluster *goRedisClusterImpl) Expire(key string, ttl time.Duration) (int64, error) {
	result, err := thisCluster.currentClusterClient.Expire(key, ttl).Result()

	if !result {
		return 0, err
	}
	return 1, nil
}

func (thisCluster *goRedisClusterImpl) Exists(key string) (bool, error) {
	result, err := thisCluster.currentClusterClient.Exists(key).Result()

	if result == 1 {
		return true, err
	}
	return false, err
}

func (thisCluster *goRedisClusterImpl) IncrXX(key string, value int64) (reply int64, err error) {
	incrExists := gocluster.NewScript(`
	if redis.call("EXISTS", KEYS[1]) ==  1 then
		return redis.call("INCRBY", KEYS[1], ARGV[1])
	else
		return -165535
	end
	`)

	result, err := incrExists.Run(thisCluster.currentClusterClient, []string{key}, value).Result()

	if reply == -165535 {
		return 0, ErrXX
	}

	return result.(int64), err
}

func (thisCluster *goRedisClusterImpl) DecrWithLimit(key string, value, lowerBound int64) (reply int64, err error) {
	decrWithLimitScript := gocluster.NewScript(`
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
	result, err := decrWithLimitScript.Run(thisCluster.currentClusterClient, []string{key}, value, lowerBound).Result()
	convRes := result.(int64)
	if err != nil {
		return lowerBound - 1, err
	}

	if convRes < lowerBound {
		return lowerBound - 1, ErrLimitExceeded
	}

	return convRes, nil
}

func (thisCluster *goRedisClusterImpl) HGetSet(key, field, value, prevValue string, ttl time.Duration) error {
	IncrByXX := gocluster.NewScript(`
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

	result, err := IncrByXX.Run(thisCluster.currentClusterClient, []string{key}, field, int(ttl), value, prevValue).Result()

	if result == valueInvalid {
		return ErrValueInvalid
	}

	return err
}

func (thisCluster *goRedisClusterImpl) ZAddToFixed(key, member string, score, maxSize int) (reply int64, err error) {
	return 0, ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) Close() error {
	return thisCluster.currentClusterClient.Close()
}

func (thisCluster *goRedisClusterImpl) Err() error {
	return ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return thisCluster.currentClusterClient.Do(commandName, args).Result()
}

func (thisCluster *goRedisClusterImpl) Send(commandName string, args ...interface{}) error {
	return ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) Flush() error {
	return ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) Receive() (reply interface{}, err error) {
	return nil, ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) GeoAdd(key string, geos ...*GeoPoint) (int64, error) {
	gclusterGeoLoc := []*gocluster.GeoLocation{}
	for _, g := range geos {
		gclusterGeoLoc = append(gclusterGeoLoc, &gocluster.GeoLocation{
			Name:      g.Member,
			Longitude: g.Longitude,
			Latitude:  g.Latitude,
		})
	}

	return thisCluster.currentClusterClient.GeoAdd(key, gclusterGeoLoc...).Result()
}

func (thisCluster *goRedisClusterImpl) GeoHash(key string, members ...string) ([]string, error) {
	return thisCluster.currentClusterClient.GeoHash(key, members...).Result()
}

func (thisCluster *goRedisClusterImpl) GeoRadius(key string, long, lat float64, q *GeoRadiusQuery) ([]*GeoLoc, error) {
	goclusterGeoRadiusQuery := &gocluster.GeoRadiusQuery{
		Radius:      q.Radius,
		Unit:        string(q.Unit),
		WithCoord:   q.WithCoord,
		WithDist:    q.WithDist,
		WithGeoHash: q.WithGeoHash,
		Count:       q.Count,
		Sort:        string(q.Sort),
	}
	result := []*GeoLoc{}
	goclusterGeoLoc, err := thisCluster.currentClusterClient.GeoRadius(key, long, lat, goclusterGeoRadiusQuery).Result()
	if err != nil {
		return result, err
	}

	for _, g := range goclusterGeoLoc {
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

func (thisCluster *goRedisClusterImpl) MGet(keys []string) ([]string, error) {
	return nil, ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) MSet(values map[string]string) error {
	return ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) MSetEx(values map[string]string, ttl time.Duration) error {
	return ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) LLen(key string) (int64, error) {
	return 0, ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) LPop(key string, count int) ([]string, error) {
	return nil, ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) LPush(key string, values []string) (int64, error) {
	return 0, ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) LPushX(key string, values []string) (int64, error) {
	return 0, ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) RPop(key string, count int) ([]string, error) {
	return nil, ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) RPush(key string, values []string) (int64, error) {
	return 0, ErrNotSupported
}

func (thisCluster *goRedisClusterImpl) RPushX(key string, values []string) (int64, error) {
	return 0, ErrNotSupported
}
