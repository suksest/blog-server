package redis

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
)

// RedisConnect connects to a default redis server at port 6379
func RedisConnect() redis.Conn {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	return c
}

func Find(t string) string {

	c := RedisConnect()
	defer c.Close()

	keys, err := c.Do("KEYS", "user:*")
	if err != nil {
		panic(err)
	}

	for _, k := range keys.([]interface{}) {
		var token string

		reply, err := c.Do("GET", k.([]byte))
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(reply.([]byte), &token); err != nil {
			panic(err)
		}

		if t == "Bearer "+token {
			return string(k.([]byte))
		}
	}
	return ""
}
