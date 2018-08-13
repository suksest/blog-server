package slidingwindowlog

import (
	"encoding/json"
	"fmt"
	"redis"
	"time"

	"github.com/labstack/echo"
)

func Init(c *Config, id string) { //id can be username for authenticated user, or IP for anonymous user
	r := redis.RedisConnect()
	defer r.Close()

	// Init
	var timelog []int64
	timelog = append(timelog, time.Now().UnixNano())
	_, err := r.Do("SET", c.Prefix+"_"+id, timelog)
	if err != nil {
		fmt.Printf(fmt.Sprint(err))
	}
}

func SetHeader(c echo.Context, limit, remain uint, reset int64) {
	c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprint(limit))
	c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprint(remain))
	c.Response().Header().Set("X-RateLimit-Update", fmt.Sprint(reset))
}

func Limiter(config *Config, c echo.Context, id string) bool {
	r := redis.RedisConnect()
	defer r.Close()

	period := GetPeriodInt(config.Period)
	reply, _ := r.Do("GET", config.Prefix+"_"+id)
	if reply == nil {
		Init(config, id)
	} else {
		var timelog []int64
		if err := json.Unmarshal(reply.([]byte), &timelog); err != nil {
			fmt.Printf(fmt.Sprint(err))
		}
		reqnum := uint(len(timelog))
		now := time.Now().UnixNano()
		for i := reqnum - 1; i >= 0; i-- { //delete expire log
			if (now - timelog[i]) > period {
				timelog = timelog[i:]
				break
			}
		}
		reqnum = uint(len(timelog)) //update request number
		if reqnum < config.MaxRequest {
			timelog = append(timelog, time.Now().UnixNano())
			SetHeader(c, config.MaxRequest, config.MaxRequest-reqnum, period-(now-timelog[0]))
			return true
		} else {
			fmt.Printf("Limit Exceeded for :" + id + "\n")
			SetHeader(c, config.MaxRequest, 0, period-now-timelog[0])
			return false
		}
	}
	SetHeader(c, config.MaxRequest, config.MaxRequest, period)
	return true
}
