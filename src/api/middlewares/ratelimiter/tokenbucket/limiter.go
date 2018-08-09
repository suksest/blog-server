package tokenbucket

import (
	"encoding/json"
	"fmt"
	"math"
	"redis"
	"time"

	"github.com/labstack/echo"
)

func Init(b *Bucket, id string) { //id can be username for authenticated user, or IP for anonymous user
	r := redis.RedisConnect()
	defer r.Close()

	// Init
	_, err := r.Do("HMSET", b.Prefix+"_"+id, "tokens", b.Capacity, "ts", fmt.Sprint(b.LastRefillTimestamp))
	if err != nil {
		panic(err)
	}
}

func SetHeader(c echo.Context, limit, remain uint) { //id can be username for authenticated user, or IP for anonymous user
	c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprint(limit))
	c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprint(remain))
}

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
			} else {
				reqTime := time.Now().Unix()
				replyToken, err := r.Do("HGET", "userlimiter_"+id, "tokens")
				if err != nil {
					panic(err)
				}
				if err := json.Unmarshal(replyToken.([]byte), &config.AvailableTokens); err != nil {
					panic(err)
				}
				replyTS, err := r.Do("HGET", "userlimiter_"+id, "ts")
				if err != nil {
					panic(err)
				}
				if err := json.Unmarshal(replyTS.([]byte), &config.LastRefillTimestamp); err != nil {
					panic(err)
				}
				// fmt.Printf(" Available tokens: " + fmt.Sprint(token) + "\n")
				if config.AvailableTokens != 0 {
					_, err := r.Do("HINCRBY", "userlimiter_"+id, "tokens", -1)
					if err != nil {
						panic(err)
					}
					config.AvailableTokens--
					_, err = r.Do("HSET", "userlimiter_"+id, "ts", reqTime)
					if err != nil {
						panic(err)
					}
					replyTS, err := r.Do("HGET", "userlimiter_"+id, "ts")
					if err != nil {
						panic(err)
					}
					if err := json.Unmarshal(replyTS.([]byte), &config.LastRefillTimestamp); err != nil {
						panic(err)
					}
					// fmt.Printf("Current tokens: " + fmt.Sprint(reply) + "\n")
					SetHeader(c, config.Capacity, config.AvailableTokens)
					return next(c)
				} else {
					elapsedTime := float64(reqTime - config.LastRefillTimestamp)
					fmt.Printf("Elapsed time (seconds): " + fmt.Sprint(elapsedTime) + "\n")
					tokensToBeAdded := elapsedTime / float64(1000) * float64(GetPeriodInt(config.Period))
					fmt.Printf("Tokens to be added: " + fmt.Sprint(tokensToBeAdded) + "\n")
					if uint(math.Floor(tokensToBeAdded)) > 0 {
						_, err := r.Do("HSET", "userlimiter_"+id, "tokens", uint(math.Floor(tokensToBeAdded)))
						if err != nil {
							panic(err)
						}
						SetHeader(c, config.Capacity, config.AvailableTokens)
						return next(c)
					} else {
						SetHeader(c, config.Capacity, 0)
						return echo.ErrForbidden
					}
				}
			}
			SetHeader(c, config.Capacity, config.AvailableTokens)
			return next(c)
		}
	}
}
