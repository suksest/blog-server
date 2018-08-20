package slidingwindowlog

import (
	"encoding/json"
	"fmt"
	"redis"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func Init(config *Config, id string) { //id can be username for authenticated user, or IP for anonymous user
	r := redis.RedisConnect()
	defer r.Close()

	// Init
	timecount := map[int64]int{
		time.Now().Unix(): 1,
	}
	timejson, _ := json.Marshal(timecount)
	fmt.Printf("timecount : %v\n", string(timejson))
	_, err := r.Do("SET", config.Prefix+"_"+id, timejson)
	if err != nil {
		fmt.Printf(fmt.Sprint(err))
	}
}

func SetHeader(c echo.Context, limit, remain int) {
	c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprint(limit))
	c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprint(remain))
}

func Limiter(config *Config, c echo.Context, id string) bool {
	r := redis.RedisConnect()
	defer r.Close()

	period := GetPeriodInt(config.Period)
	reply, _ := r.Do("GET", config.Prefix+"_"+id)
	if reply == nil {
		Init(config, id)
	} else {
		var timecount map[string]int
		if err := json.Unmarshal(reply.([]byte), &timecount); err != nil {
			fmt.Printf(fmt.Sprint(err))
		}
		// fmt.Printf("timecount: %v\n", timecount)
		now := time.Now().Unix()
		reqnum := 0
		for key, value := range timecount {
			keyint, _ := strconv.ParseInt(key, 10, 64)
			if (now - keyint) > period {
				delete(timecount, key)
			} else {
				reqnum += value
			}
		}
		if reqnum < config.MaxRequest {
			timecount[fmt.Sprint(time.Now().Unix())]++
			timejson, _ := json.Marshal(timecount)
			// fmt.Printf("timecount : %v\n", string(timejson))
			_, err := r.Do("SET", config.Prefix+"_"+id, timejson)
			if err != nil {
				fmt.Printf(fmt.Sprint(err))
			}
			SetHeader(c, config.MaxRequest, config.MaxRequest-reqnum)
			return true
		} else {
			fmt.Printf("Limit Exceeded for :" + id + "\n")
			SetHeader(c, config.MaxRequest, 0)
			return false
		}
	}
	SetHeader(c, config.MaxRequest, config.MaxRequest)
	return true
}
