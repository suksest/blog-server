package tokenbucket

import (
	"encoding/json"
	"fmt"
	"redis"
	"time"

	"github.com/labstack/echo"
)

//Init initialize hash set for user in Redis
func Init(b *Bucket, id string) { //id can be username for authenticated user, or IP for anonymous user
	r := redis.RedisConnect()
	defer r.Close()

	// Init
	_, err := r.Do("HMSET", b.Prefix+"_"+id, "tokens", b.Capacity, "ts", time.Now().Unix())
	if err != nil {
		panic(err)
	}
}

//SetHeader set header for response
func SetHeader(c echo.Context, limit, remain uint) { //id can be username for authenticated user, or IP for anonymous user
	c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprint(limit))
	c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprint(remain))
}

//Limiter limit request
func Limiter(config *Bucket, c echo.Context, id string) bool {
	r := redis.RedisConnect()
	defer r.Close()

	keys, err := r.Do("HGETALL", config.Prefix+"_"+id)
	if err != nil {
		panic(err)
	}
	if fmt.Sprint(keys) == "[]" {
		Init(config, id)
		Take(config.Prefix, id)
	} else {
		reqTime := time.Now().Unix()
		if GetTokens(config.Prefix, id) > 0 {
			Take(config.Prefix, id)
			SetHeader(c, config.Capacity, GetTokens(config.Prefix, id))
			return true
		}
		elapsedTime := GetElapsedTime(reqTime, GetLastRefillTimestamp(config.Prefix, id))
		tokensToBeAdded := GetTokensToBeAdded(elapsedTime, GetRefillInterval(config.Capacity, config.Period))
		if tokensToBeAdded > 0 {
			if tokensToBeAdded <= config.Capacity {
				Refill(config.Prefix, id, tokensToBeAdded, config.Capacity)
				Take(config.Prefix, id)
				SetHeader(c, config.Capacity, GetTokens(config.Prefix, id))
				return true
			}
			Refill(config.Prefix, id, config.Capacity, config.Capacity)
			Take(config.Prefix, id)
			SetHeader(c, config.Capacity, GetTokens(config.Prefix, id))
			return true
		}
		SetHeader(c, config.Capacity, 0)
		return false
	}
	SetHeader(c, config.Capacity, GetTokens(config.Prefix, id))
	return true
}

//GetTokens return current available token in bucket
func GetTokens(prefix, id string) uint {
	var tokens uint
	r := redis.RedisConnect()
	defer r.Close()

	reply, err := r.Do("HGET", prefix+"_"+id, "tokens")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(reply.([]byte), &tokens); err != nil {
		panic(err)
	}

	return tokens
}

//Refill bucket with token in certain period
func Refill(prefix, id string, tokens uint, capacity uint) {
	r := redis.RedisConnect()
	defer r.Close()

	currentTokens := GetTokens(prefix, id)
	if currentTokens+tokens > capacity {
		_, err := r.Do("HSET", prefix+"_"+id, "tokens", capacity)
		if err != nil {
			panic(err)
		}
	}

	_, err := r.Do("HINCRBY", prefix+"_"+id, "tokens", tokens)
	if err != nil {
		panic(err)
	}

	_, err = r.Do("HSET", prefix+"_"+id, "ts", time.Now().Unix())
	if err != nil {
		panic(err)
	}
}

//Take a token to permit a request
func Take(prefix, id string) {
	r := redis.RedisConnect()
	defer r.Close()

	_, err := r.Do("HINCRBY", prefix+"_"+id, "tokens", -1)
	if err != nil {
		panic(err)
	}
}

//GetLastRefillTimestamp return bucket's last refill timestamp
func GetLastRefillTimestamp(prefix, id string) int64 {
	var lastRefillTimestamp int64
	r := redis.RedisConnect()
	defer r.Close()

	reply, err := r.Do("HGET", prefix+"_"+id, "ts")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(reply.([]byte), &lastRefillTimestamp); err != nil {
		panic(err)
	}

	return lastRefillTimestamp
}

//GetTokensToBeAdded calculate and return number of tokens to be added
func GetTokensToBeAdded(elapsedTime int64, interval int64) uint {
	tokens := uint(elapsedTime / interval)
	return tokens
}

//GetRefillInterval calculate and return refill interval
func GetRefillInterval(capacity uint, period string) int64 {
	interval := int64(GetPeriodInt(period) / int64(capacity))
	return interval
}

//GetElapsedTime return elapsed time between bucket's last refill time and current request time
func GetElapsedTime(requestTime int64, lastRefillTimestamp int64) int64 {
	return (requestTime - lastRefillTimestamp)
}
