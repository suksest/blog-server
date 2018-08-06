package fixedWindowCounter

import (
	"fmt"
	"redis"
	"strings"
	"time"
)

func Init(t string) bool {

	r := redis.RedisConnect()
	defer r.Close()

	keys, err := r.Do("KEYS", "userlimiter:"+t+":*")
	if err != nil {
		panic(err)
	}
	if fmt.Sprint(keys) == "[]" {
		// Init
		reply, err := r.Do("SET", "userlimiter:"+t+":"+fmt.Sprint(time.Now().Unix()), 1)
		if err != nil {
			panic(err)
		}
		fmt.Printf(fmt.Sprint(reply) + "\n")
	} else {
		for _, k := range keys.([]interface{}) {
			key := string(k.([]byte))
			params := strings.Split(key, ":") //prefix:username:time

			fmt.Printf(params[2] + "\n")
			// reply, err := c.Do("GET", k.([]byte))
			// if err != nil {
			// 	panic(err)
			// }
			// if err := json.Unmarshal(reply.([]byte), &token); err != nil {
			// 	panic(err)
			// }

		}
	}

	return false
}

func UserLimiter(t string) bool {

	c := redis.RedisConnect()
	defer c.Close()

	// keys, err := c.Do("KEYS", "userlimiter:"+t+":*")
	keys, err := c.Do("KEYS", "post:*")
	if err != nil {
		panic(err)
	}

	for _, k := range keys.([]interface{}) {
		fmt.Printf(string(k.([]byte)) + "\n")

		// 	var reqtime string

		// 	reply, err := c.Do("GET", k.([]byte))
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	if err := json.Unmarshal(reply.([]byte), &token); err != nil {
		// 		panic(err)
		// 	}

		// 	if t == "Bearer "+token {
		// 		return true
		// 	}
	}
	return false
}
