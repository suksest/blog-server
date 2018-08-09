package tokenBucket

import (
	"encoding/json"
	"fmt"
	"redis"
	"time"
)

type Bucket struct {
	availableTokens int64
	startTime       time.Time
}

func Init(t string) {
	r := redis.RedisConnect()
	defer r.Close()

	// Init
	reply, err := r.Do("HMSET", "userlimiter:"+t, "tokens", 5, "ts", fmt.Sprint(time.Now().Unix()))
	if err != nil {
		panic(err)
	}

	fmt.Printf(fmt.Sprint(reply) + "\n")
}

func UserLimiter(t string) bool {
	r := redis.RedisConnect()
	defer r.Close()

	keys, err := r.Do("KEYS", "userlimiter:"+t)
	if err != nil {
		panic(err)
	}
	if fmt.Sprint(keys) == "[]" {
		Init(t)
	} else {
		var token int64
		var startTime int64
		reqTime := time.Now().Unix()
		replyToken, err := r.Do("HGET", "userlimiter:"+t, "tokens")
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(replyToken.([]byte), &token); err != nil {
			panic(err)
		}
		replyTS, err := r.Do("HGET", "userlimiter:"+t, "ts")
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(replyTS.([]byte), &startTime); err != nil {
			panic(err)
		}
		fmt.Printf(" Available tokens: " + fmt.Sprint(token) + "\n")
		if token != 0 {
			reply, err := r.Do("HINCRBY", "userlimiter:"+t, "tokens", -1)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Current tokens: " + fmt.Sprint(reply) + "\n")
			return true
		} else {
			d := (reqTime - startTime) / 60
			fmt.Printf("Elapsed time (minutes): " + fmt.Sprint(d) + "\n")
			if d < 1 {
				fmt.Printf("Token not added\n")
				return false
			} else {
				fmt.Printf("Add " + fmt.Sprint(d) + " token(s)\n")
				replyToken, err := r.Do("HSET", "userlimiter:"+t, "tokens", d)
				if err != nil {
					panic(err)
				}
				fmt.Printf("Current tokens: " + fmt.Sprint(replyToken) + "\n")
				return true
			}
		}
	}

	return false
}
