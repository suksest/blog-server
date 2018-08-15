package leakybucket

import (
	"encoding/json"
	"fmt"
	"redis"
	"time"

	"github.com/labstack/echo"
)

//Init initialize hash set for user in Redis
func Init(b *Bucket, id string) { //id can be username for authenticated user, or IP for anonymous user
	r := redis.RedisConnect()
	defer r.Close()

	// Init
	_, err := r.Do("HMSET", b.Prefix+"_"+id, "drops", 0, "allowed", 0, "ts", time.Now().Unix())
	if err != nil {
		panic(err)
	}
}

//SetHeader set header for response
func SetHeader(c echo.Context, limit, remain uint) { //id can be username for authenticated user, or IP for anonymous user
	c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprint(limit))
	c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprint(remain))
}

//Limiter limit request
func Limiter(config *Bucket, c echo.Context, id string) bool {
	r := redis.RedisConnect()
	defer r.Close()

	keys, err := r.Do("HGETALL", config.Prefix+"_"+id)
	if err != nil {
		panic(err)
	}
	if fmt.Sprint(keys) == "[]" {
		Init(config, id)
		Fill(config.Prefix, id)
		ticker := time.NewTicker(time.Duration(GetInterval(config.AllowedRequest, GetPeriodInt(config.Period))))
		go func() {
			for t := range ticker.C {
				fmt.Println(t.Unix())
				Leak(config.Prefix, id)
			}
		}()
	}
	SetHeader(c, config.Capacity, config.AllowedRequest-GetAllowed(config.Prefix, id))
	return true
}

//IsFull check the bucket is full or not
func IsFull(capacity uint, prefix, id string) bool {
	if GetDrops(prefix, id) == capacity {
		return true
	}
	return false
}

//IsEmpty check the bucket is empty or not
func IsEmpty(prefix, id string) bool {
	if GetDrops(prefix, id) == 0 {
		return true
	}
	return false
}

//IsLimitExceed check limiter condition
func IsLimitExceed(capacity, limit uint, prefix, id string) bool {
	if GetAllowed(prefix, id) == limit {
		return true
	} else if IsFull(capacity, prefix, id) {
		return true
	}
	return false
}

//Fill increment drops in bucket
func Fill(prefix, id string) {
	r := redis.RedisConnect()
	defer r.Close()

	_, err := r.Do("HINCRBY", prefix+"_"+id, "drops", 1)
	if err != nil {
		panic(err)
	}
}

//Leak decrement drops and increment allowed request in bucket
func Leak(prefix, id string) {
	r := redis.RedisConnect()
	defer r.Close()

	_, err := r.Do("HINCRBY", prefix+"_"+id, "drops", -1)
	if err != nil {
		panic(err)
	}

	_, err = r.Do("HINCRBY", prefix+"_"+id, "allowed", 1)
	if err != nil {
		panic(err)
	}
}

//UpdateTimestamp update timestamp in hash set
func UpdateTimestamp(prefix, id string, timestamp int64) {
	r := redis.RedisConnect()
	defer r.Close()

	_, err := r.Do("HSET", prefix+"_"+id, "ts", timestamp)
	if err != nil {
		panic(err)
	}
}

//GetDrops return number of drops in bucket
func GetDrops(prefix, id string) uint {
	var drops uint
	r := redis.RedisConnect()
	defer r.Close()

	reply, err := r.Do("HGET", prefix+"_"+id, "drops")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(reply.([]byte), &drops); err != nil {
		panic(err)
	}

	return drops
}

//GetTimestamp return timestamp in redis
func GetTimestamp(prefix, id string) int64 {
	var timestamp int64
	r := redis.RedisConnect()
	defer r.Close()

	reply, err := r.Do("HGET", prefix+"_"+id, "ts")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(reply.([]byte), &timestamp); err != nil {
		panic(err)
	}

	return timestamp
}

//GetAllowed return allowed request stored in redis
func GetAllowed(prefix, id string) uint {
	var allowed uint
	r := redis.RedisConnect()
	defer r.Close()

	reply, err := r.Do("HGET", prefix+"_"+id, "allowed")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(reply.([]byte), &allowed); err != nil {
		panic(err)
	}

	return allowed
}

//GetInterval return interval for each request
func GetInterval(allowed uint, duration int64) int64 {
	duration *= 1000 //duration in millisecond
	result := float64(duration) / float64(allowed)

	interval := int64(result)

	return interval
}
