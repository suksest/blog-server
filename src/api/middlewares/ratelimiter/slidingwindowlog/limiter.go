package slidingwindowlog

import (
	"encoding/json"
	"fmt"
	"redis"
	"strings"
	"time"

	"github.com/labstack/echo"
)

func Init(config *Config, id string) { //id can be username for authenticated user, or IP for anonymous user
	r := redis.RedisConnect()
	defer r.Close()

	// Init
	var timelog []int64
	timelog = append(timelog, time.Now().UnixNano())
	_, err := r.Do("SET", config.Prefix+"_"+id, timelog)
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
		var timelog []int64
		// fmt.Printf("timelog: " + string(reply.([]byte)))
		timestr := strings.Replace(string(reply.([]byte)), " ", ",", -1) //replace space to comma for valid json
		// fmt.Printf("timelog2: " + timestr)
		if err := json.Unmarshal([]byte(timestr), &timelog); err != nil {
			fmt.Printf(fmt.Sprint(err))
		}
		// fmt.Printf("timelog3: " + fmt.Sprint(timelog))
		reqnum := len(timelog)
		now := time.Now().UnixNano()
		for i := reqnum - 1; i >= 0; i-- { //delete expire log
			if (now - timelog[i]) > period {
				timelog = timelog[i:]
				break
			}
		}
		reqnum = len(timelog) //update request number
		if reqnum < config.MaxRequest {
			timelog = append(timelog, time.Now().UnixNano())
			_, err := r.Do("SET", config.Prefix+"_"+id, timelog)
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
