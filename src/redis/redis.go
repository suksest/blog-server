package redis

import (
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

// func FindAll() Posts {

// 	c := RedisConnect()
// 	defer c.Close()

// 	keys, err := c.Do("KEYS", "post:*")
// 	HandleError(err)

// 	var posts Posts

// 	for _, k := range keys.([]interface{}) {
// 		var post Post

// 		reply, err := c.Do("GET", k.([]byte))
// 		HandleError(err)
// 		if err := json.Unmarshal(reply.([]byte), &post); err != nil {
// 			panic(err)
// 		}

// 		posts = append(posts, post)
// 	}
// 	return posts
// }

// func FindPost(id int) Post {
// 	var post Post

// 	c := RedisConnect()
// 	defer c.Close()

// 	reply, err := c.Do("GET", "post:"+strconv.Itoa(id))
// 	HandleError(err)

// 	fmt.Println("GET OK")
// 	if err = json.Unmarshal(reply.([]byte), &post); err != nil {
// 		panic(err)
// 	}
// 	return post
// }
