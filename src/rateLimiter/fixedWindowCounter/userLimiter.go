package fixedWindowCounter

import (
	"encoding/json"
	"fmt"
	"redis"
	"strconv"
	"strings"
	"time"
)

var timeLimit int64 = 60
var maxCounter int = 5

func Init(t string) {
	r := redis.RedisConnect()
	defer r.Close()

	// Init
	reply, err := r.Do("SET", "userlimiter:"+t+":"+fmt.Sprint(time.Now().Unix()), 1)
	if err != nil {
		panic(err)
	}
	fmt.Printf(fmt.Sprint(reply) + "\n")
}

func UserLimiter(t string) bool {

	r := redis.RedisConnect()
	defer r.Close()

	keys, err := r.Do("KEYS", "userlimiter:"+t+":*")
	if err != nil {
		panic(err)
	}
	if fmt.Sprint(keys) == "[]" {
		Init(t)
	} else {
		for _, k := range keys.([]interface{}) {
			key := string(k.([]byte))
			params := strings.Split(key, ":") //prefix:username:time

			fmt.Printf(params[2] + "-")
			fmt.Printf(fmt.Sprint(time.Now().Unix()) + "\n")
			savedTime, err := strconv.ParseInt(params[2], 10, 64)
			if err != nil {
				panic(err)
			}
			reqTime := time.Now().Unix()
			if (reqTime - savedTime) > timeLimit {
				Init(t)
				reply, err := r.Do("DEL", k.([]byte))
				if err != nil {
					panic(err)
				}
				fmt.Printf(fmt.Sprint(reply) + "\n")
				return true
			} else {
				var counter int
				reply, err := r.Do("GET", k.([]byte))
				if err != nil {
					panic(err)
				}
				fmt.Printf(fmt.Sprint(reply) + "\n")
				if err := json.Unmarshal(reply.([]byte), &counter); err != nil {
					panic(err)
				}
				if counter < maxCounter {
					reply, err := r.Do("INCR", k.([]byte))
					if err != nil {
						panic(err)
					}
					fmt.Printf(fmt.Sprint(reply) + "\n")
					return true
				} else {
					return false
				}
			}
		}
	}
	return false
}
