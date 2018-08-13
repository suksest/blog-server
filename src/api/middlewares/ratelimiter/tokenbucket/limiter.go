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
	_, err := r.Do("HMSET", b.Prefix+"_"+id, "tokens", b.Capacity, "ts", fmt.Sprint(b.StartTimestamp))
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
func Limiter(config *Bucket) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.RealIP()
			r := redis.RedisConnect()
			defer r.Close()

			keys, err := r.Do("KEYS", config.Prefix+"_"+id)
			if err != nil {
				panic(err)
			}
			if fmt.Sprint(keys) == "[]" {
				Init(config, id)
				Take(id)
			} else {
				reqTime := time.Now().Unix()
				if GetTokens(id) != 0 {
					Take(id)
					SetHeader(c, config.Capacity, GetTokens(id))
					return next(c)
				}
				elapsedTime := GetElapsedTime(reqTime, GetLastRefillTimestamp(id))
				tokensToBeAdded := GetTokensToBeAdded(elapsedTime, config.Period)
				if tokensToBeAdded > 0 {
					Refill(id, tokensToBeAdded)
					SetHeader(c, config.Capacity, GetTokens(id))
					return next(c)
				}
				SetHeader(c, config.Capacity, 0)
				return echo.ErrForbidden
			}
			SetHeader(c, config.Capacity, GetTokens(id))
			return next(c)
		}
	}
}

//GetTokens return current available token in bucket
func GetTokens(id string) uint {
	var tokens uint
	r := redis.RedisConnect()
	defer r.Close()

	reply, err := r.Do("HGET", "userlimiter_"+id, "tokens")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(reply.([]byte), &tokens); err != nil {
		panic(err)
	}

	return tokens
}

//Refill bucket with token in certain period
func Refill(id string, tokens uint) {
	r := redis.RedisConnect()
	defer r.Close()

	_, err := r.Do("HSET", "userlimiter_"+id, "tokens", tokens)
	if err != nil {
		panic(err)
	}

	_, err = r.Do("HSET", "userlimiter_"+id, "ts", time.Now().Unix())
	if err != nil {
		panic(err)
	}
}

//Take a token to permit a request
func Take(id string) {
	r := redis.RedisConnect()
	defer r.Close()

	_, err := r.Do("HINCRBY", "userlimiter_"+id, "tokens", -1)
	if err != nil {
		panic(err)
	}
}

//GetLastRefillTimestamp return bucket's last refill timestamp
func GetLastRefillTimestamp(id string) int64 {
	var lastRefillTimestamp int64
	r := redis.RedisConnect()
	defer r.Close()

	reply, err := r.Do("HGET", "userlimiter_"+id, "ts")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(reply.([]byte), &lastRefillTimestamp); err != nil {
		panic(err)
	}

	return lastRefillTimestamp
}

//GetTokensToBeAdded calculate and return number of tokens to be added
func GetTokensToBeAdded(elapsedTime int64, period string) uint {
	tokens := uint(float64(elapsedTime) / float64(1000) * float64(GetPeriodInt(period)))
	return tokens
}

//GetElapsedTime return elapsed time between bucket's last refill time and current request time
func GetElapsedTime(requestTime int64, lastRefillTimestamp int64) int64 {
	return (requestTime - lastRefillTimestamp)
}
