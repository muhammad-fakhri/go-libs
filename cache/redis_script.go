package cache

import (
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Scripting error related
var (
	ErrClusterNotSupport = errors.New("scripting does not support on this cluster library")
	ErrLimitExceeded     = errors.New("limit exceeded")
	ErrValueInvalid      = errors.New("value invalid")
)

var incrExists = redis.NewScript(1, `
	if redis.call("EXISTS", KEYS[1]) ==  1 then
		return redis.call("INCRBY", KEYS[1], ARGV[1])
	else
		return -165535
	end
`)

var decrWithLimitScript = redis.NewScript(1, `
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

const valueInvalid = "errValueInvalid"

var hGetSetScript = redis.NewScript(3, `
	local key = KEYS[1]
	local field = KEYS[2]
	local expireSecond = tonumber(KEYS[3])
	local prevValue = redis.call('hget', key, field) or ''
	if (prevValue == ARGV[2]) then
		redis.call('hset', key, field, ARGV[1])
		redis.call('expire', key, expireSecond)
		return ARGV[1]
	else
		return 'errValueInvalid'
	end
`)

var zAddToFixed = redis.NewScript(1, `
	local key = KEYS[1]
	local score = ARGV[1]
	local member = ARGV[2]
	local max = ARGV[3]
	local maxidx = tostring(tonumber(max)-1)

	-- get smallest score
	local pmember, pscore = unpack(redis.call('zrevrange', key, maxidx, maxidx, 'withscores'))

	-- sorted set is not full
	if pmember == nil then
		redis.call('zadd', key, score, member)
		return 1

	-- sorted set is full
	else
		-- replace smallest score with new score
		if tonumber(score) > tonumber(pscore) then
			local res = redis.call('zadd', key, score, member)

			-- if res == 0 member is duplicate, else remove member with smallest score
			if res == 1 then
				redis.call('zrem', key, pmember)
			end

			return 1
		end
	end
	return 0
`)

func (r *redigoImpl) IncrXX(key string, value int64) (reply int64, err error) {
	if r.Pool == nil {
		return 0, ErrClusterNotSupport
	}

	conn := r.GetConn()
	defer conn.Close()

	reply, err = redis.Int64(incrExists.Do(conn, key, value))
	if err != nil {
		return 0, err
	}

	if reply == -165535 {
		return 0, ErrXX
	}

	return
}

func (r *redigoImpl) DecrWithLimit(key string, value, lowerBound int64) (reply int64, err error) {
	if r.Pool == nil {
		return lowerBound - 1, ErrClusterNotSupport
	}

	conn := r.GetConn()
	defer conn.Close()

	reply, err = redis.Int64(decrWithLimitScript.Do(conn, key, value, lowerBound))
	if err != nil {
		return lowerBound - 1, err
	}

	if reply < lowerBound {
		return lowerBound - 1, ErrLimitExceeded
	}

	return
}

func (r *redigoImpl) HGetSet(key, field, value, prevValue string, ttl time.Duration) error {
	if r.Pool == nil {
		return ErrClusterNotSupport
	}

	conn := r.GetConn()
	defer conn.Close()

	result, err := redis.String(hGetSetScript.Do(conn, key, field, int64(ttl.Seconds()), value, prevValue))
	if err != nil {
		return err
	}

	if result == valueInvalid {
		return ErrValueInvalid
	}

	return nil
}

func (r *redigoImpl) ZAddToFixed(key, member string, score, maxSize int) (reply int64, err error) {
	if r.Pool == nil {
		return 0, ErrClusterNotSupport
	}

	conn := r.GetConn()
	defer conn.Close()

	return redis.Int64(zAddToFixed.Do(conn, key, score, member, maxSize))
}
