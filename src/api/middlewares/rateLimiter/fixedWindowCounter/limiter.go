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
	_, err := r.Do("SET", c.Prefix+":"+id+":"+fmt.Sprint(time.Now().Unix()), 1)
	if err != nil {
		panic(err)
	}
}

func serverHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Shade/1.0")
		c.Response().Header().Set("customHeader", "thisHaveNoMeaning")

		return next(c)
	}
}

func Limiter(c *Config, id string) bool {

	r := redis.RedisConnect()
	defer r.Close()

	keys, err := r.Do("KEYS", c.Prefix+":"+id+":*")
	if err != nil {
		panic(err)
	}
	if fmt.Sprint(keys) == "[]" {
		Init(c, id)
	} else {
		for _, k := range keys.([]interface{}) {
			key := string(k.([]byte))
			params := strings.Split(key, ":") //prefix:id:time

			savedTime, err := strconv.ParseInt(params[2], 10, 64)
			if err != nil {
				panic(err)
			}
			reqTime := time.Now().Unix()
			if (reqTime - savedTime) > GetPeriodInt(c.Period) {
				Init(c, id)
				_, err := r.Do("DEL", k.([]byte))
				if err != nil {
					panic(err)
				}
				return true
			} else {
				var counter uint
				reply, err := r.Do("GET", k.([]byte))
				if err != nil {
					panic(err)
				}
				if err := json.Unmarshal(reply.([]byte), &counter); err != nil {
					panic(err)
				}
				if counter < c.MaxRequest {
					_, err := r.Do("INCR", k.([]byte))
					if err != nil {
						panic(err)
					}
					return true
				} else {
					return false
				}
			}
		}
	}
	return false
}
