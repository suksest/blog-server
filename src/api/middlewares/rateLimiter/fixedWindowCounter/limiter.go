package fixedWindowCounter

import (
	"encoding/json"
	"fmt"
	"redis"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

func Init(c *Config, id string) { //id can be username for authenticated user, or IP for anonymous user
	r := redis.RedisConnect()
	defer r.Close()

	// Init
	_, err := r.Do("SET", c.Prefix+"_"+id+"_"+fmt.Sprint(time.Now().Unix()), 1)
	if err != nil {
		fmt.Printf(fmt.Sprint(err))
	}
}

func SetHeader(c echo.Context, limit, remain uint, reset int64) {
	c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprint(limit))
	c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprint(remain))
	c.Response().Header().Set("X-RateLimit-Reset", fmt.Sprint(reset))

}

func Limiter(config *Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.RealIP()
			r := redis.RedisConnect()
			defer r.Close()

			keys, err := r.Do("KEYS", config.Prefix+"_"+id+"_*")
			if err != nil {
				fmt.Printf(fmt.Sprint(err))
			}
			if fmt.Sprint(keys) == "[]" {
				Init(config, id)
			} else {
				for _, k := range keys.([]interface{}) {
					key := string(k.([]byte))
					params := strings.Split(key, "_") //prefix_id_time

					savedTime, err := strconv.ParseInt(params[2], 10, 64)
					if err != nil {
						fmt.Printf(fmt.Sprint(err))
					}
					reqTime := time.Now().Unix()
					duration := reqTime - savedTime
					if duration > GetPeriodInt(config.Period) {
						Init(config, id)
						_, err := r.Do("DEL", k.([]byte))
						if err != nil {
							fmt.Printf(fmt.Sprint(err))
						}
					} else {
						var counter uint
						reply, err := r.Do("GET", k.([]byte))
						if err != nil {
							fmt.Printf(fmt.Sprint(err))
						}
						if err := json.Unmarshal(reply.([]byte), &counter); err != nil {
							fmt.Printf(fmt.Sprint(err))
						}
						if counter < config.MaxRequest {
							_, err := r.Do("INCR", k.([]byte))
							if err != nil {
								fmt.Printf(fmt.Sprint(err))
							}
							SetHeader(c, config.MaxRequest, config.MaxRequest-counter, GetPeriodInt(config.Period)-duration)
							return next(c)
						} else {
							fmt.Printf("Limit Exceeded for :" + id + "\n")
							SetHeader(c, config.MaxRequest, 0, GetPeriodInt(config.Period)-duration)
							return echo.ErrForbidden
						}
					}
				}
			}
			SetHeader(c, config.MaxRequest, config.MaxRequest, GetPeriodInt(config.Period))
			return next(c)
		}
	}
}
